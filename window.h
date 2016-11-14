#ifndef window_h
#define window_h

#import <Cocoa/Cocoa.h>
#import <WebKit/Webkit.h>

// This macro is used to defer the execution of a block of code in the main
// event loop.
#define defer(code)                                                            \
  dispatch_async(dispatch_get_main_queue(), ^{                                 \
                     code})

typedef struct Window__ {
  char *ID;
  char *Title;
  CGFloat X;
  CGFloat Y;
  CGFloat Width;
  CGFloat Height;
  char *BackgroundColor;
  NSVisualEffectMaterial Vibrancy;
  BOOL Borderless;
  BOOL FixedSize;
  BOOL CloseHidden;
  BOOL MinimizeHidden;
  BOOL TitlebarHidden;
  char *HTML;
  char *ResourcePath;
} Window__;

@interface WindowController
    : NSWindowController <NSWindowDelegate, WKNavigationDelegate,
                          WKScriptMessageHandler>
@property NSString *ID;
@property(weak) WKWebView *webview;
@property dispatch_semaphore_t sema;

- (instancetype)initWithID:(NSString *)ID;
@end

@interface TitleBar : NSView
@end

const void *Window_New(Window__ w);
WKWebView *Window_NewWebview(WindowController *controller, NSString *HTML,
                             NSString *resourcePath);
void Window_SetWebview(NSWindow *win, WKWebView *webview);
void Window_SetTitleBar(NSWindow *win, TitleBar *titleBar);
void Window_Mount(const void *ptr, const char *markup);
void Window_CallJS(const void *ptr, const char *js);
NSRect Window_Frame(const void *ptr);
void Window_Move(const void *ptr, CGFloat x, CGFloat y);
void Window_Resize(const void *ptr, CGFloat width, CGFloat height);
void Window_Close(const void *ptr);

#endif /* window_h */