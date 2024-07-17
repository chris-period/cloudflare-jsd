// https://github.com/daku10/go-lz-string
// Package lzstring implements the LZ-String algorithm for string compression
// and decompression. The library features two main sets of functions,
// Compress and Decompress, which are used to compress and decompress strings,
// respectively.
package cloudflare

import (
	"bytes"
	"errors"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
)

var (
	ErrInputInvalidString = errors.New("input is invalid string")
	ErrInputNotDecodable  = errors.New("input is not decodable")
	ErrInputNil           = errors.New("input should not be nil")
	ErrInputBlank         = errors.New("input should not be blank")
)

type LZString struct{}

// CompressToEncodedURIComponent takes an uncompressed string and compresses it into
// a URL-safe string, where special characters are replaced with safe alternatives.
// It returns an error if the input string is not a valid UTF-8 string.
func (lz LZString) CompressToEncodedURIComponent(uncompressed string, key string) (string, error) {
	if !utf8.ValidString(uncompressed) {
		return "", ErrInputInvalidString
	}
	res, err := lz._compress(uncompressed, 6, func(i int) []uint16 {
		return []uint16{uint16(key[i])}
	})
	if err != nil {
		return "", err
	}
	return string(utf16.Decode(res)), nil
}

type getCharFunc func(i int) []uint16

// make consistency with slice of uint16 to be enclosed with bracket.
func uint16ToString(x uint16) string {
	var b bytes.Buffer
	b.WriteByte('[')
	b.WriteString(strconv.Itoa(int(x)))
	b.WriteByte(']')
	return b.String()
}

func uint16sToString(xs []uint16) string {
	var b bytes.Buffer
	b.WriteByte('[')
	for i, x := range xs {
		b.WriteString(strconv.Itoa(int(x)))
		if i != len(xs)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteByte(']')
	return b.String()
}

func (lz LZString) _compress(uncompressed string, bitsPerChar int, getCharFromInt getCharFunc) ([]uint16, error) {
	var i, value int
	contextDictionary := make(map[string]int)
	contextDictionaryToCreate := make(map[string]bool)
	var contextC uint16
	var contextWC, contextW []uint16
	contextEnLargeIn := 2
	contextDictSize := 3
	contextNumBits := 2
	contextData := make([][]uint16, 0)
	contextDataVal := 0
	contextDataPosition := 0
	var ii int
	uncompressedRune := utf16.Encode([]rune(uncompressed))
	for ii = 0; ii < len(uncompressedRune); ii++ {
		contextC = uncompressedRune[ii]
		contextCKey := uint16ToString(contextC)
		if _, ok := contextDictionary[contextCKey]; !ok {
			contextDictionary[contextCKey] = contextDictSize
			contextDictSize++
			contextDictionaryToCreate[contextCKey] = true
		}
		contextWC = make([]uint16, len(contextW))
		copy(contextWC, contextW)
		contextWC = append(contextWC, contextC)
		contextWCKey := uint16sToString(contextWC)
		contextWKey := uint16sToString(contextW)
		if _, ok := contextDictionary[contextWCKey]; ok {
			contextW = contextWC
		} else {
			if _, ok := contextDictionaryToCreate[contextWKey]; ok {
				if len(contextW) > 0 && contextW[0] < 256 {
					for i = 0; i < contextNumBits; i++ {
						contextDataVal = contextDataVal << 1
						if contextDataPosition == bitsPerChar-1 {
							contextDataPosition = 0
							contextData = append(contextData, getCharFromInt(contextDataVal))
							contextDataVal = 0
						} else {
							contextDataPosition++
						}
					}
					value = int(contextW[0])
					for i = 0; i < 8; i++ {
						contextDataVal = (contextDataVal << 1) | (value & 1)
						if contextDataPosition == bitsPerChar-1 {
							contextDataPosition = 0
							contextData = append(contextData, getCharFromInt(contextDataVal))
							contextDataVal = 0
						} else {
							contextDataPosition++
						}
						value = value >> 1
					}
				} else {
					value = 1
					for i = 0; i < contextNumBits; i++ {
						contextDataVal = (contextDataVal << 1) | value
						if contextDataPosition == bitsPerChar-1 {
							contextDataPosition = 0
							contextData = append(contextData, getCharFromInt(contextDataVal))
							contextDataVal = 0
						} else {
							contextDataPosition++
						}
						value = 0
					}
					value = int(contextW[0])
					for i = 0; i < 16; i++ {
						contextDataVal = (contextDataVal << 1) | (value & 1)
						if contextDataPosition == bitsPerChar-1 {
							contextDataPosition = 0
							contextData = append(contextData, getCharFromInt(contextDataVal))
							contextDataVal = 0
						} else {
							contextDataPosition++
						}
						value = value >> 1
					}
				}
				contextEnLargeIn--
				if contextEnLargeIn == 0 {
					contextEnLargeIn = 1 << contextNumBits
					contextNumBits++
				}
				delete(contextDictionaryToCreate, contextWKey)
			} else {
				value = contextDictionary[contextWKey]
				for i = 0; i < contextNumBits; i++ {
					contextDataVal = (contextDataVal << 1) | (value & 1)
					if contextDataPosition == bitsPerChar-1 {
						contextDataPosition = 0
						contextData = append(contextData, getCharFromInt(contextDataVal))
						contextDataVal = 0
					} else {
						contextDataPosition++
					}
					value = value >> 1
				}
			}
			contextEnLargeIn--
			if contextEnLargeIn == 0 {
				contextEnLargeIn = 1 << contextNumBits
				contextNumBits++
			}
			contextDictionary[uint16sToString(contextWC)] = contextDictSize
			contextDictSize++
			contextW = []uint16{contextC}
		}
	}
	if len(contextW) != 0 {
		contextWKey := uint16sToString(contextW)
		if _, ok := contextDictionaryToCreate[contextWKey]; ok {
			if contextW[0] < 256 {
				for i = 0; i < contextNumBits; i++ {
					contextDataVal = contextDataVal << 1
					if contextDataPosition == bitsPerChar-1 {
						contextDataPosition = 0
						contextData = append(contextData, getCharFromInt(contextDataVal))
						contextDataVal = 0
					} else {
						contextDataPosition++
					}
				}
				value = int(contextW[0])
				for i = 0; i < 8; i++ {
					contextDataVal = (contextDataVal << 1) | (value & 1)
					if contextDataPosition == bitsPerChar-1 {
						contextDataPosition = 0
						contextData = append(contextData, getCharFromInt(contextDataVal))
						contextDataVal = 0
					} else {
						contextDataPosition++
					}
					value = value >> 1
				}
			} else {
				value = 1
				for i = 0; i < contextNumBits; i++ {
					contextDataVal = (contextDataVal << 1) | value
					if contextDataPosition == bitsPerChar-1 {
						contextDataPosition = 0
						contextData = append(contextData, getCharFromInt(contextDataVal))
						contextDataVal = 0
					} else {
						contextDataPosition++
					}
					value = 0
				}
				value = int(contextW[0])
				for i = 0; i < 16; i++ {
					contextDataVal = (contextDataVal << 1) | (value & 1)
					if contextDataPosition == bitsPerChar-1 {
						contextDataPosition = 0
						contextData = append(contextData, getCharFromInt(contextDataVal))
						contextDataVal = 0
					} else {
						contextDataPosition++
					}
					value = value >> 1
				}
			}
			contextEnLargeIn--
			if contextEnLargeIn == 0 {
				contextEnLargeIn = 1 << contextNumBits
				contextNumBits++
			}
			delete(contextDictionaryToCreate, contextWKey)
		} else {
			value = contextDictionary[contextWKey]
			for i = 0; i < contextNumBits; i++ {
				contextDataVal = (contextDataVal << 1) | (value & 1)
				if contextDataPosition == bitsPerChar-1 {
					contextDataPosition = 0
					contextData = append(contextData, getCharFromInt(contextDataVal))
					contextDataVal = 0
				} else {
					contextDataPosition++
				}
				value = value >> 1
			}
		}
		contextEnLargeIn--
		if contextEnLargeIn == 0 {
			// original algorithm has below expression, but this value is unused probably.
			// contextEnLargeIn = 1 << contextNumBits
			contextNumBits++
		}
	}

	value = 2
	for i = 0; i < contextNumBits; i++ {
		contextDataVal = (contextDataVal << 1) | (value & 1)
		if contextDataPosition == bitsPerChar-1 {
			contextDataPosition = 0
			contextData = append(contextData, getCharFromInt(contextDataVal))
			contextDataVal = 0
		} else {
			contextDataPosition++
		}
		value = value >> 1
	}

	for {
		contextDataVal = contextDataVal << 1
		if contextDataPosition == bitsPerChar-1 {
			contextData = append(contextData, getCharFromInt(contextDataVal))
			break
		} else {
			contextDataPosition++
		}
	}
	result := make([]uint16, 0)
	for _, cd := range contextData {
		result = append(result, cd...)
	}
	return result, nil
}
