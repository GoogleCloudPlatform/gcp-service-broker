package common

type Kind int

const (
	NorthAmericaNorthEast1 Kind = iota
	UsCentral
	UsCentral1
	UsEast1
	UsEast4
	UsWest1
	UsWest2
	SouthAmericaEast1
	EuropeNorth1
	EuropeWest1
	EuropeWest2
	EuropeWest3
	EuropeWest4
	EuropeWest6
	AsiaEast1
	AsiaEast2
	AsiaNorthEast1
	AsiaNorthEast2
	AsiaSouth1
	AsiaSouthEast1
	AustraliaSouthEast1
)

func (k Kind) Region() string {
	switch k {
	case NorthAmericaNorthEast1:
		return "northamerica-northeast1"
	case UsCentral:
		return "us-central"
	case UsCentral1:
		return "us-central1"
	case UsEast1:
		return "us-east1"
	case UsEast4:
		return "us-east4"
	case UsWest1:
		return "us-west1"
	case UsWest2:
		return "us-west2"
	case SouthAmericaEast1:
		return "southamerica-east1"
	case EuropeNorth1:
		return "europe-north1"
	case EuropeWest1:
		return "europe-west1"
	case EuropeWest2:
		return "europe-west2"
	case EuropeWest3:
		return "europe-west3"
	case EuropeWest4:
		return "europe-west4"
	case EuropeWest6:
		return "europe-west6"
	case AsiaEast1:
		return "asia-east1"
	case AsiaEast2:
		return "asia-east2"
	case AsiaNorthEast1:
		return "asia-northeast1"
	case AsiaNorthEast2:
		return "asia-northeast2"
	case AsiaSouth1:
		return "asia-south1"
	case AsiaSouthEast1:
		return "asia-southeast1"
	case AustraliaSouthEast1:
		return "australia-southeast1"
	default:
		return ""
	}
}
