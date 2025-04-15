package jobs

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/samber/lo"

	"github.com/BabyLev/Umka-1/internal/repo/satellites"
)

// Задача: раз в сутки запрашивать информацию обо всех спутниках в хранилище
// у r4uab сервера и обновлять встроенный (embedded) объект Satellite

func (j *Jobs) UpdateSatellitesInfo(ctx context.Context) {
	// шаг 1. Получаем все спутники из хранилища +
	sats, err := j.repoSats.FindSatellite(ctx, satellites.FilterSatellite{
		NoradIDNotNull: lo.ToPtr(true),
	})
	if err != nil {
		log.Default().Printf("j.repoSats.FindSatellite: %s", err.Error())
		return
	}
	// начало цикла
	for _, sat := range sats {
		if sat.NoradID == nil {
			continue
		}

		reqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		updatedSatInfo, err := j.r4uabClient.GetSatelliteInfo(reqCtx, *sat.NoradID)
		if err != nil {
			log.Default().Printf("j.r4uabClient.GetSatelliteInfo: %s", err.Error())
		}

		noradID, err := strconv.ParseInt(updatedSatInfo.SatelliteId, 10, 64)
		if err != nil {
			log.Default().Printf("strconv.ParseInt: %s", err.Error())
		}

		err = j.repoSats.UpdateSatellite(ctx, satellites.Satellite{
			ID:      sat.ID,
			SatName: updatedSatInfo.Name,
			NoradID: &noradID,
			Line1:   sat.Line1,
			Line2:   sat.Line2,
		})
		if err != nil {
			log.Default().Printf("j.storage.UpdateSatellite: %s", err.Error())
		}
	}

	// шаг 2(ПРОВЕРКА) Если Norad ID != nil, переходим к шагу 3, иначе к следующему спутнику в цикле
	// шаг 3. Получаем спутник с новыми значениями из r4uab по Norad ID взятого спутника
	// шаг 4. Создаем новый объект Satellite (из пакета Satellite) на основе данных, взятых из r4uab
	// шаг 5. Обновляем запись в хранилище
}
