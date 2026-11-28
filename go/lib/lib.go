package lib

type Fun[A, R any] = func(...A) (R, error)
