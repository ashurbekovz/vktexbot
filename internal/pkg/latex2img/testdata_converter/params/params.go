package params

import "github.com/shopspring/decimal"

func ImageDPI() decimal.Decimal {
    return decimal.NewFromInt(400)
}

func CorrectTestdataFiles() []string {
    return []string{
        "simple.tex",
        "double_compilation.tex",
    }
}


