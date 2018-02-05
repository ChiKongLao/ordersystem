package datamodels

// 邮件
const (
	From = "From"
	To = "To"
	Cc = "Cc"
	Subject = "Subject"
)

type Mail struct {
	From string			// 发件人
	To []string			// 收件人
	Cc []string			// 抄送人
	Subject string		// 主题
	Body string			// 正文
	Attach string		// 附件
}
