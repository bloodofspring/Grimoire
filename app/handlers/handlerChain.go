package handlers

import (
	"context"
	"log"
	"time"

	tele "gopkg.in/telebot.v4"
)

type Arg map[string]any

type chainHandler func(c tele.Context, args *Arg) (*Arg, error)

type execLog struct {
	Success bool
	Message string
	Error   error
}

type HandlerChain struct {
	Handlers      []chainHandler
	Args          *Arg
	ExecutionLogs execLog
	timeout       time.Duration
}

func (hc HandlerChain) Init(timeout time.Duration, handlers ...chainHandler) *HandlerChain {
	log.Printf("HandlerChain.Init: Creating HandlerChain with timeout: %v", timeout)

	new := HandlerChain{
		Handlers: handlers,
		timeout:  timeout,
		Args:     &Arg{},
	}

	log.Printf("HandlerChain.Init: HandlerChain created with %d handlers", len(handlers))
	return &new
}

func (hc *HandlerChain) Run(c tele.Context) error {
	log.Printf("HandlerChain.Run: Starting execution with timeout: %v", hc.timeout)

	// Создаем контекст в момент выполнения
	ctx, cancel := context.WithTimeout(context.Background(), hc.timeout)
	defer cancel()

	deadline, ok := ctx.Deadline()
	log.Printf("HandlerChain.Run: Context created, deadline: %v", deadline)
	log.Printf("HandlerChain.Run: Context deadline ok: %v", ok)

	done := make(chan error, 1)

	go func() {
		log.Printf("HandlerChain.Run: Starting goroutine with %d handlers", len(hc.Handlers))
		for i, handler := range hc.Handlers {
			log.Printf("HandlerChain.Run: Executing handler %d", i+1)

			// Проверяем контекст перед выполнением каждого хендлера
			select {
			case <-ctx.Done():
				log.Printf("HandlerChain.Run: Context cancelled before handler %d", i+1)
				hc.ExecutionLogs.Success = false
				hc.ExecutionLogs.Message = "Context cancelled"
				hc.ExecutionLogs.Error = ctx.Err()
				done <- ctx.Err()
				return
			default:
				log.Printf("HandlerChain.Run: Context is still active before handler %d", i+1)
			}

			newArgs, err := handler(c, hc.Args)
			if err != nil {
				log.Printf("HandlerChain.Run: Handler %d failed: %v", i+1, err)
				hc.ExecutionLogs.Success = false
				hc.ExecutionLogs.Message = err.Error()
				hc.ExecutionLogs.Error = err
				done <- err
				return
			}
			*hc.Args = *newArgs
			log.Printf("HandlerChain.Run: Handler %d completed successfully", i+1)
		}
		log.Printf("HandlerChain.Run: All handlers completed successfully")
		done <- nil
	}()

	// Ждем либо завершения выполнения, либо отмены контекста
	select {
	case err := <-done:
		log.Printf("HandlerChain.Run: Goroutine completed with result: %v", err)
		return err
	case <-ctx.Done():
		log.Printf("HandlerChain.Run: Context timeout occurred")
		// Контекст отменен по таймауту
		hc.ExecutionLogs.Success = false
		hc.ExecutionLogs.Message = "Context timeout"
		hc.ExecutionLogs.Error = ctx.Err()
		return ctx.Err()
	}
}
