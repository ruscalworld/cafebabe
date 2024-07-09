package cafebabe

import (
	"io"
	"reflect"
)

type ClassFile struct {
	Magic        U4
	MinorVersion U2
	MajorVersion U2

	ConstantPoolCount U2
	ConstantPool      []CpInfo `cb:"1-indexed"`

	AccessFlags U2
	ThisClass   U2
	SuperClass  U2

	//InterfacesCount U2
	//Interfaces      []U2
	//
	//FieldsCount U2
	//Fields      []FieldInfo
	//
	//MethodsCount U2
	//Methods      []MethodInfo
	//
	//AttributesCount U2
	//Attributes      []AttributeInfo
}

type CpInfoTag U1

func (c *CpInfoTag) ReadFrom(r io.Reader) (n int64, err error) {
	u1 := U1(*c)

	_, err = u1.ReadFrom(r)
	if err != nil {
		return n, err
	}

	*c = CpInfoTag(u1)
	return n, nil
}

const (
	ConstantClass              CpInfoTag = 7
	ConstantFieldRef           CpInfoTag = 9
	ConstantMethodRef          CpInfoTag = 10
	ConstantInterfaceMethodRef CpInfoTag = 11
	ConstantString             CpInfoTag = 8
	ConstantInteger            CpInfoTag = 3
	ConstantFloat              CpInfoTag = 4
	ConstantLong               CpInfoTag = 5
	ConstantDouble             CpInfoTag = 6
	ConstantNameAndType        CpInfoTag = 12
	ConstantUtf8               CpInfoTag = 1
	ConstantMethodHandle       CpInfoTag = 15
	ConstantMethodType         CpInfoTag = 16
	ConstantInvokeDynamic      CpInfoTag = 18
)

type CpInfo struct {
	Tag  CpInfoTag
	Info any `cb:"variadic"`
}

func (i *CpInfo) PrepareInfo() reflect.Type {
	switch i.Tag {
	case ConstantClass:
		return reflect.TypeFor[ConstantClassInfo]()
	case ConstantMethodRef, ConstantFieldRef, ConstantInterfaceMethodRef:
		return reflect.TypeFor[ConstantRefInfo]()
	case ConstantString:
		return reflect.TypeFor[ConstantStringInfo]()
	case ConstantInteger:
		return reflect.TypeFor[ConstantIntegerInfo]()
	case ConstantFloat:
		return reflect.TypeFor[ConstantFloatInfo]()
	case ConstantLong:
		return reflect.TypeFor[ConstantLongInfo]()
	case ConstantDouble:
		return reflect.TypeFor[ConstantDoubleInfo]()
	case ConstantNameAndType:
		return reflect.TypeFor[ConstantNameAndTypeInfo]()
	case ConstantUtf8:
		return reflect.TypeFor[ConstantUtf8Info]()
	case ConstantMethodHandle:
		return reflect.TypeFor[ConstantMethodHandleInfo]()
	case ConstantMethodType:
		return reflect.TypeFor[ConstantMethodTypeInfo]()
	case ConstantInvokeDynamic:
		return reflect.TypeFor[ConstantInvokeDynamicInfo]()
	}

	return nil
}

type ConstantClassInfo struct {
	// Index of UTF-8 encoded class name in constant pool
	NameIndex U2
}

type ConstantRefInfo struct {
	// Index of ClassInfo in constant pool
	ClassIndex U2

	// Index of NameAndType in constant pool
	NameAndTypeIndex U2
}

type ConstantStringInfo struct {
	// Index of UTF-8 encoded string in constant pool
	StringIndex U2
}

type ConstantIntegerInfo struct {
	// Big-endian ordered bytes of an integer
	Bytes U4
}

type ConstantFloatInfo struct {
	// Big-endian ordered IEEE 754 floating-point number
	Bytes U4
}

type ConstantLongInfo struct {
	HighBytes U4
	LowBytes  U4
}

type ConstantDoubleInfo struct {
	HighBytes U4
	LowBytes  U4
}

type ConstantNameAndTypeInfo struct {
	// Index of UTF-8 encoded name in constant pool
	NameIndex U2

	// Index of UTF-8 encoded descriptor in constant pool
	DescriptorIndex U2
}

type ConstantUtf8Info struct {
	Length U2
	Bytes  []U1
}

func (i *ConstantUtf8Info) String() string {
	return string(i.Bytes)
}

type ConstantMethodHandleInfo struct {
	ReferenceKind U1

	// Index of reference in constant pool
	ReferenceIndex U2
}

type ConstantMethodTypeInfo struct {
	// Index of UTF-8 encoded descriptor in constant pool
	DescriptorIndex U2
}

type ConstantInvokeDynamicInfo struct {
	// Index of bootstrap method in BootstrapMethods table
	BootstrapMethodAttrIndex U2

	// Index of NameAndType in constant pool
	NameAndTypeIndex U2
}

type FieldInfo struct {
	AccessFlags     U2
	NameIndex       U2
	DescriptorIndex U2

	AttributesCount U2
	Attributes      []AttributeInfo
}

type MethodInfo struct {
}

type AttributeInfo struct {
}
