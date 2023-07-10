package util

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertToProtoTimestamp(input time.Time) *timestamppb.Timestamp {
	return timestamppb.New(input)
}

func ConvertInt32SliceToIntSlice(input []int32) []int {
	output := make([]int, len(input))

	for index, val := range input {
		output[index] = int(val)
	}

	return output
}

func ConvertIntSliceToInt32Slice(input []int) []int32 {
	output := make([]int32, len(input))

	for index, val := range input {
		output[index] = int32(val)
	}

	return output
}

func ConvertInt32MapToIntMap(input map[int32]int32) map[int]int {
	output := make(map[int]int)

	for key, value := range input {
		output[int(key)] = int(value)
	}

	return output
}

func ConvertIntMapToInt32Map(input map[int]int) map[int32]int32 {
	output := make(map[int32]int32)

	for key, value := range input {
		output[int32(key)] = int32(value)
	}

	return output
}
