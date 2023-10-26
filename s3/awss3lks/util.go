package awss3lks

import "strings"

const (
	alphabet            = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-._"
	replacementAlphabet = "-_0"
)

func wellFormBlobName(n string) (string, error) {
	if n == "" {
		return n, nil
	}

	var out strings.Builder
	var err error
	lenReplacementAlphabet := int32(len(replacementAlphabet))
	for _, c := range n {
		if strings.Index(alphabet, string(c)) >= 0 {
			_, err = out.WriteString(string(c))
		} else {
			if c > lenReplacementAlphabet {
				c = c % lenReplacementAlphabet
			} else {
				c = lenReplacementAlphabet % c
			}
			err = out.WriteByte(replacementAlphabet[c])
		}

		if err != nil {
			return "", err
		}
	}

	return out.String(), nil
}
