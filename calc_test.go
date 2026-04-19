
package main

import "testing"

func TestDivide(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        res, err := Divide(10, 2)
        if err != nil || res != 5 {
            t.Errorf("expected 5, got %d", res)
        }
    })

    t.Run("divide by zero", func(t *testing.T) {
        _, err := Divide(10, 0)
        if err == nil {
            t.Errorf("expected error")
        }
    })
}

func TestSubtractTable(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 5, 3, 2},
        {"zero", 5, 0, 5},
        {"neg+pos", -1, 4, -5},
        {"both neg", -5, -2, -3},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Subtract(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("got %d want %d", got, tt.want)
            }
        })
    }
}
