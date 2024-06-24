package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/sign"
)

type shaWriter struct {
	w          http.ResponseWriter
	b          *bytes.Buffer
	key        string
	statusCode int
}

func newSignatureWriter(w http.ResponseWriter, k string) *shaWriter {
	var b bytes.Buffer
	return &shaWriter{
		w:          w,
		key:        k,
		b:          &b,
		statusCode: http.StatusOK,
	}
}

func (s *shaWriter) Header() http.Header {
	return s.w.Header()
}
func (s *shaWriter) Write(p []byte) (int, error) {
	n, err := s.b.Write(p)

	if err != nil {
		return n, fmt.Errorf("shaWriter write error: %w", err)
	}
	return n, nil
}
func (s *shaWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode
}

// Close добавляет заголовок-подпись и пишет данные в ответ из буфера.
func (s *shaWriter) Close() error {
	signature, err := sign.Sign(s.b.Bytes(), s.key)
	if err != nil {
		return fmt.Errorf("signing error: %w", err)
	}
	s.w.Header().Set("Hashsha256", signature)
	s.w.WriteHeader(s.statusCode)
	_, err = s.w.Write(s.b.Bytes())
	if err != nil {
		return fmt.Errorf("shaWriter error writing response: %w", err)
	}
	return nil
}

func shaMiddlewareBuilder(key string) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hashHeader := r.Header.Get("Hash")
			supportsSigning := (hashHeader != "none" && hashHeader != "")
			sw := w

			if supportsSigning {
				// проверяем подпись
				sentSignature := r.Header.Get("Hashsha256")
				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Info(fmt.Errorf("error reading body: %w", err))
					return
				}
				err = r.Body.Close()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Info(fmt.Errorf("error closing body: %w", err))
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(body))
				expectedSignature, err := sign.Sign(body, key)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Info(fmt.Errorf("error signing: %w", err))
					return
				}
				if sentSignature != expectedSignature {
					w.WriteHeader(http.StatusBadRequest)
					_, err = w.Write([]byte("wrong signature"))
					if err != nil {
						logger.Info(fmt.Errorf("error writing response: %w", err))
					}
					return
				}

				// подписываем тело при закрытии
				sw := newSignatureWriter(w, key)
				defer func() {
					_ = sw.Close()
				}()
			}
			next.ServeHTTP(sw, r)
		})
	}
}
