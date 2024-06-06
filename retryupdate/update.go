//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var oldValue *string = nil // старое значение
	var authErr *kvapi.AuthError
	var conflictErr *kvapi.ConflictError
	var oldVersion uuid.UUID
	for finishGet := false; !finishGet; {
		get, err := c.Get(&kvapi.GetRequest{Key: key})
		switch {
		//если ключ не найден все плохо
		case errors.Is(err, kvapi.ErrKeyNotFound):
			finishGet = true

		//если ошибка доступа
		case errors.As(err, &authErr):
			return err

			// если все ок
		case err == nil:
			oldValue = &get.Value
			oldVersion = get.Version
			finishGet = true
		}
	}
	newValue, err := updateFn(oldValue)
	if err != nil {
		return err
	}
	newVersion := uuid.Must(uuid.NewV4())
	for finishSet := false; !finishSet; {
		_, err := c.Set(&kvapi.SetRequest{
			Key:        key,
			Value:      newValue,
			OldVersion: oldVersion,
			NewVersion: newVersion,
		})
		switch {
		case errors.As(err, &conflictErr):
			// значение изменилось с момента последнего чтения
			if conflictErr.ExpectedVersion != newVersion {
				return UpdateValue(c, key, updateFn)
			} else {
				return nil
			}

		case errors.As(err, &authErr):
			return err
		case err == nil:
			return nil
		case errors.Is(err, kvapi.ErrKeyNotFound):
			oldVersion = uuid.UUID{}
			newValue, err = updateFn(nil)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
