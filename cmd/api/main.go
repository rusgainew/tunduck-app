package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main - точка входа в приложение
func main() {
	// Создаем контекст корневого уровня
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем и инициализируем приложение с контекстом
	app, err := NewApp(ctx, ".env")
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	// Канал для приема сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Горутина для запуска сервера
	errChan := make(chan error, 1)
	go func() {
		errChan <- app.Run()
	}()

	// Ожидаем либо ошибку сервера, либо сигнал завершения
	select {
	case <-sigChan:
		log.Println("Получен сигнал завершения, инициируем graceful shutdown...")
		// Отменяем контекст для начала процесса завершения
		cancel()
	case err := <-errChan:
		// Обрабатываем ошибку сервера
		if err != nil {
			log.Printf("Ошибка запуска сервера: %v", err)
			log.Println("Инициируем graceful shutdown...")
		}
		// Отменяем контекст для начала процесса завершения
		cancel()
	}

	// Выполняем graceful shutdown с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("Ошибка при завершении сервера: %v", err)
	}

	log.Println("Приложение корректно завершило работу")
}
