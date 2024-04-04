package utils

var (
	IdVerify               = Rules{"ID": []string{NotEmpty()}}
	ApiVerify              = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}
	MenuVerify             = Rules{"Path": {NotEmpty()}, "ParentId": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	MenuMetaVerify         = Rules{"Title": {NotEmpty()}}
	LoginVerify            = Rules{"CaptchaId": {}, "Username": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify         = Rules{"Username": {NotEmpty()}, "NickName": {NotEmpty()}, "Password": {NotEmpty()}, "AuthorityId": {NotEmpty()}}
	PageInfoVerify         = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	CustomerVerify         = Rules{"CustomerName": {NotEmpty()}, "CustomerPhoneData": {NotEmpty()}}
	AutoCodeVerify         = Rules{"Abbreviation": {NotEmpty()}, "StructName": {NotEmpty()}, "PackageName": {NotEmpty()}, "Fields": {NotEmpty()}}
	AutoPackageVerify      = Rules{"PackageName": {NotEmpty()}}
	AuthorityVerify        = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}}
	AuthorityIdVerify      = Rules{"AuthorityId": {NotEmpty()}}
	OldAuthorityVerify     = Rules{"OldAuthorityId": {NotEmpty()}}
	ChangePasswordVerify   = Rules{"Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	SetUserAuthorityVerify = Rules{"AuthorityId": {NotEmpty()}}
	GetUserTicketsVerify   = Rules{"category": {NotEmpty()}}
	ModelFieldVerify       = Rules{"Flag": {NotEmpty()}}
	GetModelListVerify     = Rules{"Page": {NotEmpty()}, "LIMIT": {NotEmpty()}}
	UserRegisterVerify     = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}, "Email": {NotEmpty(), RegexpMatch("^([a-zA-Z0-9]+)@[a-zA-Z0-9]+\\.[a-zA-Z]{3}$")}}
	UploadModelStoreVerify = Rules{"ModelName": {NotEmpty()}, "UUID": {NotEmpty()}, "ModelVersion": {NotEmpty()}, "ModelDescription": {NotEmpty()},
		"ModelType": {NotEmpty()}, "IsImage": {NotEmpty()}, "Cmd": {NotEmpty()}, "JsonUrl": {NotEmpty()}, "ImgUrl": {NotEmpty()}, "TaskID": {NotEmpty()}, "AlgorithmID": {NotEmpty()}}
)
