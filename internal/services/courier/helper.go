package courier

import mc "book-nest/internal/models/courier"

func GetInstantCourierOnly(input []mc.Courier) (result []mc.Courier) {
	for _, v := range input {
		if v.CourierServiceName == "Instant" {
			result = append(result, v)
		}
	}
	return result
}
