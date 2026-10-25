package main

// Тонкая обёртка над функциями Windows (только стандартная библиотека Go).

import (
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	advapi32 = syscall.NewLazyDLL("advapi32.dll")
	dwmapi   = syscall.NewLazyDLL("dwmapi.dll")
	uxtheme  = syscall.NewLazyDLL("uxtheme.dll")
	comctl32 = syscall.NewLazyDLL("comctl32.dll")

	pRegisterClassExW       = user32.NewProc("RegisterClassExW")
	pCreateWindowExW        = user32.NewProc("CreateWindowExW")
	pDefWindowProcW         = user32.NewProc("DefWindowProcW")
	pGetMessageW            = user32.NewProc("GetMessageW")
	pTranslateMessage       = user32.NewProc("TranslateMessage")
	pDispatchMessageW       = user32.NewProc("DispatchMessageW")
	pPostQuitMessage        = user32.NewProc("PostQuitMessage")
	pPostMessageW           = user32.NewProc("PostMessageW")
	pSendMessageW           = user32.NewProc("SendMessageW")
	pShowWindow             = user32.NewProc("ShowWindow")
	pMoveWindow             = user32.NewProc("MoveWindow")
	pSetWindowPos           = user32.NewProc("SetWindowPos")
	pGetClientRect          = user32.NewProc("GetClientRect")
	pInvalidateRect         = user32.NewProc("InvalidateRect")
	pBeginPaint             = user32.NewProc("BeginPaint")
	pEndPaint               = user32.NewProc("EndPaint")
	pFillRect               = user32.NewProc("FillRect")
	pDrawTextW              = user32.NewProc("DrawTextW")
	pGetDC                  = user32.NewProc("GetDC")
	pReleaseDC              = user32.NewProc("ReleaseDC")
	pSetFocus               = user32.NewProc("SetFocus")
	pGetFocus               = user32.NewProc("GetFocus")
	pIsWindowVisible        = user32.NewProc("IsWindowVisible")
	pIsWindow               = user32.NewProc("IsWindow")
	pIsIconic               = user32.NewProc("IsIconic")
	pGetForegroundWindow    = user32.NewProc("GetForegroundWindow")
	pSetForegroundWindow    = user32.NewProc("SetForegroundWindow")
	pRegisterHotKey         = user32.NewProc("RegisterHotKey")
	pUnregisterHotKey       = user32.NewProc("UnregisterHotKey")
	pSetTimer               = user32.NewProc("SetTimer")
	pKillTimer              = user32.NewProc("KillTimer")
	pCreatePopupMenu        = user32.NewProc("CreatePopupMenu")
	pAppendMenuW            = user32.NewProc("AppendMenuW")
	pTrackPopupMenu         = user32.NewProc("TrackPopupMenu")
	pDestroyMenu            = user32.NewProc("DestroyMenu")
	pGetCursorPos           = user32.NewProc("GetCursorPos")
	pScreenToClient         = user32.NewProc("ScreenToClient")
	pClientToScreen         = user32.NewProc("ClientToScreen")
	pMessageBoxW            = user32.NewProc("MessageBoxW")
	pLoadCursorW            = user32.NewProc("LoadCursorW")
	pFindWindowW            = user32.NewProc("FindWindowW")
	pGetKeyState            = user32.NewProc("GetKeyState")
	pGetWindowTextLengthW   = user32.NewProc("GetWindowTextLengthW")
	pGetWindowTextW         = user32.NewProc("GetWindowTextW")
	pSetWindowTextW         = user32.NewProc("SetWindowTextW")
	pRegisterWindowMessageW = user32.NewProc("RegisterWindowMessageW")
	pOpenClipboard          = user32.NewProc("OpenClipboard")
	pCloseClipboard         = user32.NewProc("CloseClipboard")
	pEmptyClipboard         = user32.NewProc("EmptyClipboard")
	pSetClipboardData       = user32.NewProc("SetClipboardData")
	pGetWindowPlacement     = user32.NewProc("GetWindowPlacement")
	pSetWindowPlacement     = user32.NewProc("SetWindowPlacement")
	pGetSystemMetrics       = user32.NewProc("GetSystemMetrics")
	pSetRect                = user32.NewProc("SetRect")
	pDestroyWindow          = user32.NewProc("DestroyWindow")

	pCreateFontW      = gdi32.NewProc("CreateFontW")
	pCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	pCreatePen        = gdi32.NewProc("CreatePen")
	pSelectObject     = gdi32.NewProc("SelectObject")
	pSetTextColor     = gdi32.NewProc("SetTextColor")
	pSetBkColor       = gdi32.NewProc("SetBkColor")
	pSetBkMode        = gdi32.NewProc("SetBkMode")
	pRoundRect        = gdi32.NewProc("RoundRect")
	pGetDeviceCaps    = gdi32.NewProc("GetDeviceCaps")
	pDeleteObject     = gdi32.NewProc("DeleteObject")
	pEllipse          = gdi32.NewProc("Ellipse")
	pMoveToEx         = gdi32.NewProc("MoveToEx")
	pLineTo           = gdi32.NewProc("LineTo")
	pGetStockObject   = gdi32.NewProc("GetStockObject")

	pGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
	pLoadLibraryW     = kernel32.NewProc("LoadLibraryW")
	pCreateMutexW     = kernel32.NewProc("CreateMutexW")
	pGlobalAlloc      = kernel32.NewProc("GlobalAlloc")
	pGlobalLock       = kernel32.NewProc("GlobalLock")
	pGlobalUnlock     = kernel32.NewProc("GlobalUnlock")
	pGetLastError     = kernel32.NewProc("GetLastError")
	pCloseHandle      = kernel32.NewProc("CloseHandle")

	pShellNotifyIconW = shell32.NewProc("Shell_NotifyIconW")
	pShellExecuteW    = shell32.NewProc("ShellExecuteW")
	pExtractIconExW   = shell32.NewProc("ExtractIconExW")

	pRegOpenKeyExW    = advapi32.NewProc("RegOpenKeyExW")
	pRegSetValueExW   = advapi32.NewProc("RegSetValueExW")
	pRegDeleteValueW  = advapi32.NewProc("RegDeleteValueW")
	pRegQueryValueExW = advapi32.NewProc("RegQueryValueExW")
	pRegCloseKey      = advapi32.NewProc("RegCloseKey")

	pDwmSetWindowAttribute = dwmapi.NewProc("DwmSetWindowAttribute")
	pSetWindowTheme        = uxtheme.NewProc("SetWindowTheme")
	pInitCommonControlsEx  = comctl32.NewProc("InitCommonControlsEx")
)

const (
	WS_OVERLAPPEDWINDOW = 0x00CF0000
	WS_CHILD            = 0x40000000
	WS_VISIBLE          = 0x10000000
	WS_VSCROLL          = 0x00200000
	WS_CLIPCHILDREN     = 0x02000000
	WS_TABSTOP          = 0x00010000

	ES_AUTOHSCROLL = 0x0080
	ES_MULTILINE   = 0x0004
	ES_READONLY    = 0x0800
	ES_NOHIDESEL   = 0x0100
	ES_AUTOVSCROLL = 0x0040

	LBS_NOTIFY            = 0x0001
	LBS_OWNERDRAWVARIABLE = 0x0020
	LBS_NOINTEGRALHEIGHT  = 0x0100

	WM_DESTROY         = 0x0002
	WM_SIZE            = 0x0005
	WM_SETFOCUS        = 0x0007
	WM_SETREDRAW       = 0x000B
	WM_PAINT           = 0x000F
	WM_CLOSE           = 0x0010
	WM_ERASEBKGND      = 0x0014
	WM_GETMINMAXINFO   = 0x0024
	WM_DRAWITEM        = 0x002B
	WM_MEASUREITEM     = 0x002C
	WM_SETFONT         = 0x0030
	WM_NOTIFY          = 0x004E
	WM_CONTEXTMENU     = 0x007B
	WM_KEYDOWN         = 0x0100
	WM_SYSKEYDOWN      = 0x0104
	WM_COMMAND         = 0x0111
	WM_TIMER           = 0x0113
	WM_VSCROLL         = 0x0115
	WM_CTLCOLOREDIT    = 0x0133
	WM_CTLCOLORLISTBOX = 0x0134
	WM_CTLCOLORSTATIC  = 0x0138
	WM_LBUTTONUP       = 0x0202
	WM_LBUTTONDBLCLK   = 0x0203
	WM_RBUTTONUP       = 0x0205
	WM_HOTKEY          = 0x0312
	WM_APP             = 0x8000
	WM_USER            = 0x0400

	EM_SETSEL              = 0x00B1
	EM_GETSEL              = 0x00B0
	EM_SCROLLCARET         = 0x00B7
	EM_LINESCROLL          = 0x00B6
	EM_GETFIRSTVISIBLELINE = 0x00CE
	EM_LINEINDEX           = 0x00BB
	EM_SETRECT             = 0x00B3
	EM_SETCUEBANNER        = 0x1501
	EM_EXLIMITTEXT         = WM_USER + 53
	EM_EXLINEFROMCHAR      = WM_USER + 54
	EM_SETBKGNDCOLOR       = WM_USER + 67
	EM_SETEVENTMASK        = WM_USER + 69
	EM_GETTEXTRANGE        = WM_USER + 75
	EM_AUTOURLDETECT       = WM_USER + 91
	EM_SETTEXTEX           = WM_USER + 97
	EM_SETZOOM             = WM_USER + 225
	EM_CHARFROMPOS         = 0x00D7
	EM_EXSETSEL            = WM_USER + 55
	EM_EXGETSEL            = WM_USER + 52

	EN_CHANGE    = 0x0300
	EN_MSGFILTER = 0x0700
	EN_LINK      = 0x070B

	ENM_KEYEVENTS   = 0x00010000
	ENM_MOUSEEVENTS = 0x00020000
	ENM_LINK        = 0x04000000

	LB_ADDSTRING    = 0x0180
	LB_RESETCONTENT = 0x0184
	LB_SETCURSEL    = 0x0186
	LB_GETCURSEL    = 0x0188
	LB_GETITEMDATA  = 0x0199
	LB_SETTOPINDEX  = 0x0197
	LBN_SELCHANGE   = 1

	SW_HIDE       = 0
	SW_SHOW       = 5
	SW_RESTORE    = 9
	SW_SHOWNORMAL = 1

	SWP_NOSIZE     = 0x0001
	SWP_NOMOVE     = 0x0002
	SWP_NOACTIVATE = 0x0010
	SWP_SHOWWINDOW = 0x0040

	MOD_ALT      = 0x0001
	MOD_CONTROL  = 0x0002
	MOD_SHIFT    = 0x0004
	MOD_WIN      = 0x0008
	MOD_NOREPEAT = 0x4000

	VK_BACK      = 0x08
	VK_RETURN    = 0x0D
	VK_ESCAPE    = 0x1B
	VK_PRIOR     = 0x21
	VK_NEXT      = 0x22
	VK_END       = 0x23
	VK_HOME      = 0x24
	VK_UP        = 0x26
	VK_DOWN      = 0x28
	VK_F5        = 0x74
	VK_CONTROL   = 0x11
	VK_OEM_PLUS  = 0xBB
	VK_OEM_MINUS = 0xBD
	VK_ADD       = 0x6B
	VK_SUBTRACT  = 0x6D

	SB_LINEUP   = 0
	SB_LINEDOWN = 1
	SB_PAGEUP   = 2
	SB_PAGEDOWN = 3
	SB_TOP      = 6
	SB_BOTTOM   = 7

	DT_LEFT         = 0x0000
	DT_RIGHT        = 0x0002
	DT_VCENTER      = 0x0004
	DT_WORDBREAK    = 0x0010
	DT_SINGLELINE   = 0x0020
	DT_NOPREFIX     = 0x0800
	DT_CALCRECT     = 0x0400
	DT_END_ELLIPSIS = 0x8000

	TRANSPARENT = 1

	NIM_ADD     = 0
	NIM_MODIFY  = 1
	NIM_DELETE  = 2
	NIF_MESSAGE = 0x01
	NIF_ICON    = 0x02
	NIF_TIP     = 0x04
	NIF_INFO    = 0x10
	NIIF_INFO   = 0x01

	MF_STRING       = 0x0000
	MF_SEPARATOR    = 0x0800
	MF_CHECKED      = 0x0008
	MF_GRAYED       = 0x0001
	TPM_RETURNCMD   = 0x0100
	TPM_RIGHTBUTTON = 0x0002

	MB_OK          = 0x0000
	MB_ICONWARNING = 0x0030
	MB_ICONINFO    = 0x0040
	MB_TOPMOST     = 0x40000

	CF_UNICODETEXT = 13
	GMEM_MOVEABLE  = 0x0002

	HKEY_CURRENT_USER = 0x80000001
	KEY_QUERY_VALUE   = 0x0001
	KEY_SET_VALUE     = 0x0002
	REG_SZ            = 1

	IDC_ARROW         = 32512
	LOGPIXELSY        = 90
	FW_NORMAL         = 400
	FW_SEMIBOLD       = 600
	FW_BOLD           = 700
	DEFAULT_CHARSET   = 1
	CLEARTYPE_QUALITY = 5

	ODS_SELECTED         = 0x0001
	CS_DBLCLKS           = 0x0008
	ERROR_ALREADY_EXISTS = 183
)

var HWND_TOPMOST = ^uintptr(0)   // -1
var HWND_NOTOPMOST = ^uintptr(1) // -2

type POINT struct{ X, Y int32 }
type RECT struct{ Left, Top, Right, Bottom int32 }

type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
	Private uint32
}

type PAINTSTRUCT struct {
	Hdc         uintptr
	FErase      int32
	RcPaint     RECT
	FRestore    int32
	FIncUpdate  int32
	RgbReserved [32]byte
}

type DRAWITEMSTRUCT struct {
	CtlType    uint32
	CtlID      uint32
	ItemID     uint32
	ItemAction uint32
	ItemState  uint32
	HwndItem   uintptr
	HDC        uintptr
	RcItem     RECT
	ItemData   uintptr
}

type MEASUREITEMSTRUCT struct {
	CtlType    uint32
	CtlID      uint32
	ItemID     uint32
	ItemWidth  uint32
	ItemHeight uint32
	ItemData   uintptr
}

type MINMAXINFO struct {
	PtReserved     POINT
	PtMaxSize      POINT
	PtMaxPosition  POINT
	PtMinTrackSize POINT
	PtMaxTrackSize POINT
}

type NMHDR struct {
	HwndFrom uintptr
	IdFrom   uintptr
	Code     uint32
}

type CHARRANGE struct{ CpMin, CpMax int32 }

type ENLINK struct {
	Nmhdr  NMHDR
	Msg    uint32
	WParam uintptr
	LParam uintptr
	Chrg   CHARRANGE
}

type MSGFILTER struct {
	Nmhdr  NMHDR
	Msg    uint32
	WParam uintptr
	LParam uintptr
}

type TEXTRANGEW struct {
	Chrg      CHARRANGE
	LpstrText *uint16
}

type SETTEXTEX struct {
	Flags    uint32
	Codepage uint32
}

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

type NOTIFYICONDATAW struct {
	CbSize           uint32
	HWnd             uintptr
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            uintptr
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GuidItem         GUID
	HBalloonIcon     uintptr
}

type WINDOWPLACEMENT struct {
	Length           uint32
	Flags            uint32
	ShowCmd          uint32
	PtMinPosition    POINT
	PtMaxPosition    POINT
	RcNormalPosition RECT
}

type INITCOMMONCONTROLSEX struct {
	DwSize uint32
	DwICC  uint32
}

func u16(s string) *uint16  { p, _ := syscall.UTF16PtrFromString(s); return p }
func u16p(s string) uintptr { return uintptr(unsafe.Pointer(u16(s))) }

func SendMessage(h uintptr, msg uint32, w, l uintptr) uintptr {
	r, _, _ := pSendMessageW.Call(h, uintptr(msg), w, l)
	return r
}
func PostMessage(h uintptr, msg uint32, w, l uintptr) {
	pPostMessageW.Call(h, uintptr(msg), w, l)
}
func RGB(r, g, b int) uintptr { return uintptr(r | g<<8 | b<<16) }

func copyUTF16(dst []uint16, s string) {
	src, _ := syscall.UTF16FromString(s)
	n := copy(dst, src)
	if n == len(dst) {
		dst[len(dst)-1] = 0
	}
}

func getWindowText(h uintptr) string {
	n, _, _ := pGetWindowTextLengthW.Call(h)
	buf := make([]uint16, n+1)
	pGetWindowTextW.Call(h, uintptr(unsafe.Pointer(&buf[0])), n+1)
	return syscall.UTF16ToString(buf)
}

func messageBox(h uintptr, text, title string, flags uintptr) {
	pMessageBoxW.Call(h, u16p(text), u16p(title), flags)
}

func setClipboard(h uintptr, s string) bool {
	data, _ := syscall.UTF16FromString(s)
	if r, _, _ := pOpenClipboard.Call(h); r == 0 {
		return false
	}
	defer pCloseClipboard.Call()
	pEmptyClipboard.Call()
	size := uintptr(len(data) * 2)
	mem, _, _ := pGlobalAlloc.Call(GMEM_MOVEABLE, size)
	if mem == 0 {
		return false
	}
	ptr, _, _ := pGlobalLock.Call(mem)
	dst := unsafe.Slice((*uint16)(unsafe.Pointer(ptr)), len(data))
	copy(dst, data)
	pGlobalUnlock.Call(mem)
	r, _, _ := pSetClipboardData.Call(CF_UNICODETEXT, mem)
	return r != 0
}

func createFont(heightPx int, weight int, face string) uintptr {
	f, _, _ := pCreateFontW.Call(uintptr(int32(-heightPx)), 0, 0, 0, uintptr(weight), 0, 0, 0,
		DEFAULT_CHARSET, 0, 0, CLEARTYPE_QUALITY, 0, u16p(face))
	return f
}

func drawText(hdc uintptr, s string, rc *RECT, flags uintptr) {
	p, _ := syscall.UTF16FromString(s)
	pDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&p[0])), uintptr(len(p)-1), uintptr(unsafe.Pointer(rc)), flags)
}

func fillRect(hdc uintptr, rc *RECT, brush uintptr) {
	pFillRect.Call(hdc, uintptr(unsafe.Pointer(rc)), brush)
}

// ---- реестр (автозапуск) ----
const runKey = `Software\Microsoft\Windows\CurrentVersion\Run`

func regOpenRun(access uintptr) uintptr {
	var h uintptr
	r, _, _ := pRegOpenKeyExW.Call(HKEY_CURRENT_USER, u16p(runKey), 0, access, uintptr(unsafe.Pointer(&h)))
	if r != 0 {
		return 0
	}
	return h
}

func autostartEnabled(name string) bool {
	h := regOpenRun(KEY_QUERY_VALUE)
	if h == 0 {
		return false
	}
	defer pRegCloseKey.Call(h)
	r, _, _ := pRegQueryValueExW.Call(h, u16p(name), 0, 0, 0, 0)
	return r == 0
}

func setAutostart(name, cmd string, on bool) bool {
	h := regOpenRun(KEY_SET_VALUE)
	if h == 0 {
		return false
	}
	defer pRegCloseKey.Call(h)
	if !on {
		r, _, _ := pRegDeleteValueW.Call(h, u16p(name))
		return r == 0
	}
	data, _ := syscall.UTF16FromString(cmd)
	r, _, _ := pRegSetValueExW.Call(h, u16p(name), 0, REG_SZ, uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)*2))
	return r == 0
}
