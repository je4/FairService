package zsearch

import "fmt"

func ItemTypeFromString(s string) (ItemType, error) {
	if it, ok := ItemTypeReverse[s]; ok {
		return it, nil
	} else {
		return ItemTypeOther, fmt.Errorf("unknown item type %s", s)
	}
}

func PersonRoleFromString(s string) (PersonRole, error) {
	if it, ok := CreatorTypeReverse[s]; ok {
		return it, nil
	} else {
		return PersonRoleArtist, fmt.Errorf("unknown creator type %s", s)
	}
}
