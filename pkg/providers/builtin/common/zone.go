package common

const (
	UsCentral1A Kind = iota
	UsCentral1B
	UsCentral1C
	UsCentral1F
	UsWest2A
	UsWest2B
	UsWest2C
	UsEast4A
	UsEast4B
	UsEast4C
	UsWest1A
	UsWest1B
	UsWest1C
	UsEast1B
	UsEast1C
	UsEast1D
	NorthamericaNorthEast1A
	NorthamericaNorthEast1B
	NorthamericaNorthEast1C
	SouthAmericaEast1A
	SouthAmericaEast1B
	SouthAmericaEast1C
	EuropeWest1B
	EuropeWest1D
	EuropeNorth1A
	EuropeNorth1B
	EuropeNorth1C
	EuropeWest2A
	EuropeWest2B
	EuropeWest2C
	EuropeWest4A
	EuropeWest4B
	EuropeWest4C
	EuropeWest6A
	EuropeWest6B
	EuropeWest6C
	AsiaSouth1A
	AsiaSouth1B
	AsiaSouth1C
	AsiaSouthEast1A
	AsiaSouthEast1B
	AsiaSouthEast1C
	AsiaEast1A
	AsiaEast1B
	AsiaEast1C
	AsiaEast2A
	AsiaEast2B
	AsiaEast2C
	AsiaNorthEast1A
	AsiaNorthEast1B
	AsiaNorthEast1C
	AsiaNorthEast2A
	AsiaNorthEast2B
	AsiaNorthEast2C
	AustraliaSouthEast1A
	AustraliaSouthEast1B
	AustraliaSouthEast1C
)

func (k Kind) Zone() string {
	switch k {
	case UsCentral1A:
		return "us-central1-a"
	case UsCentral1B:
		return "us-central1-b"
	case UsCentral1C:
		return "us-central1-c"
	case UsCentral1F:
		return "us-central1-f"
	case UsWest2A:
		return "us-west2-a"
	case UsWest2B:
		return "us-west2-b"
	case UsWest2C:
		return "us-west2-c"
	case UsEast4A:
		return "us-east4-a"
	case UsEast4B:
		return "us-east4-b"
	case UsEast4C:
		return "us-east4-c"
	case UsWest1A:
		return "us-west1-a"
	case UsWest1B:
		return "us-west1-b"
	case UsWest1C:
		return "us-west1-c"
	case UsEast1B:
		return "us-east1-b"
	case UsEast1C:
		return "us-east1-c"
	case UsEast1D:
		return "us-east1-d"
	case NorthamericaNorthEast1A:
		return "northamerica-northeast1-a"
	case NorthamericaNorthEast1B:
		return "northamerica-northeast1-b"
	case NorthamericaNorthEast1C:
		return "northamerica-northeast1-c"
	case SouthAmericaEast1A:
		return "southamerica-east1-a"
	case SouthAmericaEast1B:
		return "southamerica-east1-b"
	case SouthAmericaEast1C:
		return "southamerica-east1-c"
	case EuropeWest1B:
		return "europe-west1-b"
	case EuropeWest1D:
		return "europe-west1-d"
	case EuropeNorth1A:
		return "europe-north1-a"
	case EuropeNorth1B:
		return "europe-north1-b"
	case EuropeNorth1C:
		return "europe-north1-c"
	case EuropeWest2A:
		return "europe-west2-a"
	case EuropeWest2B:
		return "europe-west2-b"
	case EuropeWest2C:
		return "europe-west2-c"
	case EuropeWest4A:
		return "europe-west4-a"
	case EuropeWest4B:
		return "europe-west4-b"
	case EuropeWest4C:
		return "europe-west4-c"
	case EuropeWest6A:
		return "europe-west6-a"
	case EuropeWest6B:
		return "europe-west6-b"
	case EuropeWest6C:
		return "europe-west6-c"
	case AsiaSouth1A:
		return "asia-south1-a"
	case AsiaSouth1B:
		return "asia-south1-b"
	case AsiaSouth1C:
		return "asia-south1-c"
	case AsiaSouthEast1A:
		return "asia-southeast1-a"
	case AsiaSouthEast1B:
		return "asia-southeast1-b"
	case AsiaSouthEast1C:
		return "asia-southeast1-c"
	case AsiaEast1A:
		return "asia-east1-a"
	case AsiaEast1B:
		return "asia-east1-b"
	case AsiaEast1C:
		return "asia-east1-c"
	case AsiaEast2A:
		return "asia-east2-a"
	case AsiaEast2B:
		return "asia-east2-b"
	case AsiaEast2C:
		return "asia-east2-c"
	case AsiaNorthEast1A:
		return "asia-northeast1-a"
	case AsiaNorthEast1B:
		return "asia-northeast1-b"
	case AsiaNorthEast1C:
		return "asia-northeast1-c"
	case AsiaNorthEast2A:
		return "asia-northeast2-a"
	case AsiaNorthEast2B:
		return "asia-northeast2-b"
	case AsiaNorthEast2C:
		return "asia-northeast2-c"
	case AustraliaSouthEast1A:
		return "australia-southeast1-a"
	case AustraliaSouthEast1B:
		return "australia-southeast1-b"
	case AustraliaSouthEast1C:
		return "australia-southeast1-c"
	default:
		return ""
	}
}
