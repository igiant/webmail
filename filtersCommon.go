package webmail

// List of all possible targets whose state can be

// FilterConditionType - tested as filter's initial condition.
type FilterConditionType string

const (
	CtEnvelopeRecipient FilterConditionType = "CtEnvelopeRecipient" // Recipient from SMTP envelope 'RCPT TO:'
	CtEnvelopeSender    FilterConditionType = "CtEnvelopeSender"    // Sender from SMTP envelope 'MAIL FROM:'
	CtRecipient         FilterConditionType = "CtRecipient"         // content of message headers 'To' and 'Cc'
	CtSender            FilterConditionType = "CtSender"            // content of message header 'Sender'
	CtFrom              FilterConditionType = "CtFrom"              // content of message header 'From'
	CtCc                FilterConditionType = "CtCc"                // content of message header 'Cc'
	CtTo                FilterConditionType = "CtTo"                // content of message header 'To'
	CtSubject           FilterConditionType = "CtSubject"           // messages with certain subject
	CtAttachment        FilterConditionType = "CtAttachment"        // messages that have attachment
	CtSize              FilterConditionType = "CtSize"              // messages with certain size (in Bytes)
	CtSpam              FilterConditionType = "CtSpam"              // messages marked as spam
	CtAll               FilterConditionType = "CtAll"               // all messages
)

// List of all possible comparators of targets (see ConditionTargetType)

// FilterComparatorType - which can be used to test filter's initial conditions.
type FilterComparatorType string

const (
	CcEqual        FilterComparatorType = "CcEqual"        // tests whether target (ex. sender) is EQUAL to some string
	CcContain      FilterComparatorType = "CcContain"      // tests whether target (ex. sender) CONTAINS to some string
	CcNotContain   FilterComparatorType = "CcNotContain"   // tests whether target (ex. sender) NOT CONTAINS to some string
	CcNotEqual     FilterComparatorType = "CcNotEqual"     // tests whether target (ex. sender) IS NOT to some string
	CcUnder        FilterComparatorType = "CcUnder"        // tests whether size of mail is UNDER some number of Bytes
	CcOver         FilterComparatorType = "CcOver"         // tests whether size of mail is OVER some number of Bytes
	CcNoComparator FilterComparatorType = "CcNoComparator" // marks conditions that do not need any comparator (ex. has attachment, for all messages and is spam)
)

// Defines initial condition of filter rule. Filter's initial
// conditions are tested whenever we need find out whether

// FilterCondition - filter should be applied or not.
type FilterCondition struct {
	TestedTarget FilterConditionType  `json:"testedTarget"` // some aspect of message that will be tested using selected comparator and against parameters
	Comparator   FilterComparatorType `json:"comparator"`   // type of target comparator (has to be set correctly, even for CtAttachment and CtAll - CcNoComparator)
	Parameters   StringList           `json:"parameters"`   // zero, one or more values that will be used to compare to target using selected comparator
}

// All possible types of filter's actions. Type of action effects number of parameters
// that this action requires. Every action has sequence of parameters. Types defined
// by FilterActionType enumeration hold information about required parameters as comment.
// @Warning: According to the selected type of action, you have to provide (set) all parameters

// FilterActionType - to filter action that this type of action requires.
type FilterActionType string

const (
	FaAddHeader     FilterActionType = "FaAddHeader"     // two mandatory parameters: 1. header name, 2. header value (Headers 'Content-*' are forbidden)
	FaSetHeader     FilterActionType = "FaSetHeader"     // two mandatory parameters: 1. header name, 2. header value (Headers 'Content-*' are forbidden)
	FaRemoveHeader  FilterActionType = "FaRemoveHeader"  // one mandatory parameter: 1. header name (Headers 'Content-*' and 'Receive' are forbidden)
	FaAddRecipient  FilterActionType = "FaAddRecipient"  // one mandatory parameter: 1. address
	FaCopyToAddress FilterActionType = "FaCopyToAddress" // one mandatory parameter: 1. address
	FaReject        FilterActionType = "FaReject"        // one mandatory parameter: 1. reason (for mail rejection)
	FaFileInto      FilterActionType = "FaFileInto"      // one mandatory parameter: 1. path to directory where to store mail
	FaRedirect      FilterActionType = "FaRedirect"      // one mandatory parameter: 1. address where to redirect mail
	FaDiscard       FilterActionType = "FaDiscard"       // no parameters required
	FaKeep          FilterActionType = "FaKeep"          // no parameters required
	FaNotify        FilterActionType = "FaNotify"        // three mandatory parameters: 1. address; 2. subject; 3. text of notification
	FaSetReadFlag   FilterActionType = "FaSetReadFlag"   // no parameters required
	FaAutoReply     FilterActionType = "FaAutoReply"     // one mandatory parameter: 1. text of message
	FaStop          FilterActionType = "FaStop"          // no parameters required
)

// Defines action of filter rule. Filter's action are
// performed when all filter's initial conditions are

// FilterAction - meet.
type FilterAction struct {
	Type       FilterActionType `json:"type"`
	Parameters StringList       `json:"parameters"` // list of parameters (see FilterActionType for more info)
}

// EvaluationModeType - Types of filter's conditions evaluation.
type EvaluationModeType string

const (
	EmAnyOf EvaluationModeType = "EmAnyOf"
	EmAllOf EvaluationModeType = "EmAllOf"
)

type FilterActionList []FilterAction

type FilterConditionList []FilterCondition
