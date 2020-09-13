package lib

// ----------------------------------------------------------------
func (this *Mlrval) GetType() MVType {
	return this.mvtype
}

func (this *Mlrval) GetTypeName() string {
	return TYPE_NAMES[this.mvtype]
}

func GetTypeName(mvtype MVType) string {
	return TYPE_NAMES[mvtype]
}

// ----------------------------------------------------------------
func (this *Mlrval) IsError() bool {
	return this.mvtype == MT_ERROR
}

func (this *Mlrval) IsAbsent() bool {
	return this.mvtype == MT_ABSENT
}

func (this *Mlrval) IsVoid() bool {
	return this.mvtype == MT_VOID
}

func (this *Mlrval) IsErrorOrVoid() bool {
	return this.mvtype == MT_ERROR || this.mvtype == MT_VOID
}

func (this *Mlrval) IsString() bool {
	return this.mvtype == MT_STRING
}

func (this *Mlrval) IsInt() bool {
	return this.mvtype == MT_INT
}

func (this *Mlrval) IsFloat() bool {
	return this.mvtype == MT_FLOAT
}

func (this *Mlrval) IsBool() bool {
	return this.mvtype == MT_BOOL
}

func (this *Mlrval) IsTrue() bool {
	return this.mvtype == MT_BOOL && this.boolval == true
}
func (this *Mlrval) IsFalse() bool {
	return this.mvtype == MT_BOOL && this.boolval == false
}

func (this *Mlrval) IsArray() bool {
	return this.mvtype == MT_ARRAY
}
func (this *Mlrval) IsMap() bool {
	return this.mvtype == MT_MAP
}
func (this *Mlrval) IsArrayOrMap() bool {
	return this.mvtype == MT_ARRAY || this.mvtype == MT_MAP
}

// ----------------------------------------------------------------
func (this *Mlrval) GetIntValue() (intValue int64, isInt bool) {
	if this.mvtype == MT_INT {
		return this.intval, true
	} else {
		return -999, false
	}
}

func (this *Mlrval) GetBoolValue() (boolValue bool, isBoolean bool) {
	if this.mvtype == MT_BOOL {
		return this.boolval, true
	} else {
		return false, false
	}
}

func (this *Mlrval) GetArray() []Mlrval {
	if this.mvtype == MT_ARRAY {
		return this.arrayval
	} else {
		return nil
	}
}

func (this *Mlrval) GetMap() *Mlrmap {
	if this.mvtype == MT_MAP {
		return this.mapval
	} else {
		return nil
	}
}

// ----------------------------------------------------------------
func (this *Mlrval) Copy() *Mlrval {
	that := *this
	if this.mvtype == MT_MAP {
		that.mapval = this.mapval.Copy()
	} else if this.mvtype == MT_ARRAY {
		that.arrayval = CopyMlrvalArray(this.arrayval)
	}
	return &that
}

func CopyMlrvalArray(input []Mlrval) []Mlrval {
	output := make([]Mlrval, len(input))
	for i, element := range input {
		output[i] = *element.Copy()
	}
	return output
}