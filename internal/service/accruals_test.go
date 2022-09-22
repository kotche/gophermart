package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccrualOrderService_CheckLuhn(t *testing.T) {

	var repo AccrualOrderRepoContract
	a := &AccrualOrderService{
		repo: repo,
	}

	tests := []struct {
		name   string
		number uint64
		want   bool
	}{
		{
			name:   "correct_number",
			number: 4561261212345467,
			want:   true,
		},
		{
			name:   "incorrect number",
			number: 1111,
			want:   false,
		},
		{
			name:   "empty number",
			number: 0,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, a.CheckLuhn(tt.number), "CheckLuhn(%v)", tt.number)
		})
	}
}
