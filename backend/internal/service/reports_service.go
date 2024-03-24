package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/persistence"
	logger "github.com/SOAT1StackGoLang/Hackaton/pkg/middleware"
	"time"
)

type rS struct {
	tR persistence.TimekeepingRepository
}

func (r rS) GetReportCSVByRangeAndUserID(ctx context.Context, userID string, start, end time.Time) ([]byte, error) {
	var (
		outBytes  []byte
		outBuffer = new(bytes.Buffer)

		rTKRStruct = &models.RangedTimekeepingReport{
			UserID:        userID,
			Start:         start,
			End:           end,
			WorkedMinutes: 0,
			Open:          false,
		}

		details = make([]models.Timekeeping, 0)
	)

	regs, err := r.tR.ListTimekeepingByRangeAndUserID(ctx, userID, start, end)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed getting timekeeping", err.Error()))
		return nil, err
	}

	for _, reg := range regs {
		if reg.Open {
			rTKRStruct.Open = true
		}
		details = append(details, *reg)
		rTKRStruct.WorkedMinutes += reg.WorkedMinutes
	}
	rTKRStruct.Details = details

	err = writeReportCSV(outBuffer, rTKRStruct) // Escreva diretamente no buffer
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing report to bytes", err.Error()))
		return nil, err
	}

	outBytes = outBuffer.Bytes()
	fmt.Println(string(outBytes))

	return outBytes, nil
}

func writeReportCSV(outBuffer *bytes.Buffer, tkrStruct *models.RangedTimekeepingReport) error {
	writer := csv.NewWriter(outBuffer)

	header := []string{"Usuário", "Começo", "Fim", "Tempo trabalhado", "Status"}
	if err := writer.Write(header); err != nil {
		return err
	}

	row := []string{
		tkrStruct.UserID,
		tkrStruct.Start.In(location).Format(time.RFC822),
		tkrStruct.End.In(location).Format(time.RFC822),
		formatarDiferencaTempo(tkrStruct.WorkedMinutes),
		parseOpen(tkrStruct.Open),
	}
	err := writer.Write(row)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		return err
	}

	// Linha em branco
	err = writer.Write([]string{""})
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		return err
	}

	err = writer.Write([]string{"Registros diários"})
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		return err
	}

	for _, detail := range tkrStruct.Details {
		err = writeTimekeepingCSV(writer, &detail)
		if err != nil {
			return err
		}
	}

	writer.Flush() // Chame Flush após escrever no buffer
	return nil
}

func writeRangedTimekeepingReportCSV(writer *csv.Writer, tkrStruct *models.RangedTimekeepingReport) error {
	header := []string{"Usuário", "Começo", "Fim", "Tempo trabalhado", "Status"}
	if err := writer.Write(header); err != nil {
		return err
	}

	row := []string{
		tkrStruct.UserID,
		tkrStruct.Start.In(location).String(),
		tkrStruct.End.In(location).String(),
		formatarDiferencaTempo(tkrStruct.WorkedMinutes),
		parseOpen(tkrStruct.Open),
	}
	err := writer.Write(row)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		return err
	}

	// Linha em branco
	err = writer.Write([]string{""})
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		return err
	}

	err = writer.Write([]string{"Registros diários"})
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		return err
	}

	for _, detail := range tkrStruct.Details {
		err = writeTimekeepingCSV(writer, &detail)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeTimekeepingCSV(writer *csv.Writer, t *models.Timekeeping) error {
	header := []string{"ID", "Data de Criação", "Data de referência", "Atualizado", "Tempo trabalhado", "Status"}

	if err := writer.Write(header); err != nil {
		return err
	}

	row := []string{
		t.ID.String(),
		t.CreatedAt.In(location).Format(time.RFC822),
		t.ReferenceDate.Format("2006-01-02"),
		t.UpdatedAt.In(location).Format(time.RFC822),
		formatarDiferencaTempo(t.WorkedMinutes),
		parseOpen(t.Open),
	}
	err := writer.Write(row)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		return err
	}

	err = writer.Write([]string{""})
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		return err
	}
	for _, detail := range t.Details {
		err = writeDetailsCSV(writer, detail)
		if err != nil {
			return err
		}

	}

	// Linha vazia após detalhe
	err = writer.Write([]string{""})
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		return err
	}

	return err
}

func writeDetailsCSV(writer *csv.Writer, detail *models.Details) error {
	header := []string{"Tempo trabalhado:", formatarDiferencaTempo(detail.WorkedMinutes)}
	err := writer.Write(header)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		return err
	}

	err = writer.Write([]string{"Entrada inicial", detail.StartingEntry.CreatedAt.In(location).Format(time.RFC822)})
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		return err
	}
	if detail.EndingEntry != nil {
		err = writer.Write([]string{"Entrada final", detail.EndingEntry.CreatedAt.In(location).Format(time.RFC)})
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %s", "failed writing row", err.Error()))
		}
	}

	return err
}

func parseOpen(in bool) string {
	if in {
		return "Há registro(s) aberto(s)"
	}
	return "Registro(s) concluído(s)"
}

func (r rS) GetReportByReferenceDateAndUserID(ctx context.Context, userID string, referenceDate time.Time) (*models.Timekeeping, error) {
	out, err := r.tR.GetTimekeepingByReferenceDateAndUserID(ctx, userID, referenceDate)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed getting timekeeping", err.Error()))
		return nil, err
	}
	return out, nil
}

func (r rS) GetReportByRangeAndUserID(ctx context.Context, userID string, start, end time.Time) (*models.RangedTimekeepingReport, error) {
	var (
		out = &models.RangedTimekeepingReport{
			UserID:        userID,
			Start:         start,
			End:           end,
			WorkedMinutes: 0,
			Open:          false,
		}

		details = make([]models.Timekeeping, 0)
	)

	regs, err := r.tR.ListTimekeepingByRangeAndUserID(ctx, userID, start, end)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s", "failed getting timekeeping", err.Error()))
		return nil, err
	}

	for _, reg := range regs {
		if reg.Open {
			out.Open = true
		}
		details = append(details, *reg)
		out.WorkedMinutes += reg.WorkedMinutes
	}
	out.Details = details

	return out, nil
}

func NewReportService(tR persistence.TimekeepingRepository) ReportService {
	return &rS{tR: tR}
}

func formatarDiferencaTempo(minutos int64) string {
	// Extrair o número de horas e minutos
	horas := minutos / 60
	minutosRestantes := minutos % 60

	// Criar a string formatada
	stringFormatada := ""
	if horas > 0 {
		stringFormatada = fmt.Sprintf("%d hora(s)", horas)
		if minutosRestantes > 0 {
			stringFormatada += " e "
		}
	}
	if minutosRestantes > 0 {
		stringFormatada += fmt.Sprintf("%d minuto(s)", minutosRestantes)
	}

	return stringFormatada
}
