package jobs

import (
	"context"
	"log"
	"time"

	"github.com/BabyLev/Umka-1/internal/types"
	"github.com/BabyLev/Umka-1/satellite"
)

// Задача: раз в сутки запрашивать информацию обо всех спутниках в хранилище
// у r4uab сервера и обновлять встроенный (embedded) объект Satellite

func (j *Jobs) UpdateSatellitesInfo() {
	// шаг 1. Получаем все спутники из хранилища +
	sats := j.storage.FindSatellite("")
	// начало цикла
	for storageSatID, sat := range sats {
		if sat.NoradID == nil {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		updatedSatInfo, err := j.r4uabClient.GetSatelliteInfo(ctx, *sat.NoradID)
		if err != nil {
			log.Default().Printf("j.r4uabClient.GetSatelliteInfo: %s", err.Error())
		}
		newSat := satellite.New(updatedSatInfo.Line1, updatedSatInfo.Line2)
		err = j.storage.UpdateSatellite(storageSatID, types.Satellite{
			Satellite: newSat,
			Name:      sat.Name,
			NoradID:   sat.NoradID,
		})
		if err != nil {
			log.Default().Printf("j.storage.UpdateSatellite: %s", err.Error())
		}
	}

	// шаг 2(ПРОВЕРКА) Если Norad ID != nil, переходим к шагу 3, иначе к следующему спутнику в цикле
	// шаг 3. Получаем спутник с новыми значениями из r4uab по Norad ID взятого спутника
	// шаг 4. Создаем новый объект Satellite (из пакета Satellite) на основе данных, взятых из r4uab
	// шаг 5. Обновляем запись в хранилище

	// что потребуется
	//
}
