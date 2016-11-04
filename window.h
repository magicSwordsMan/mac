#ifndef window_h
#define window_h

#import <Cocoa/Cocoa.h>
#import <WebKit/Webkit.h>

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
} Window__;

@interface WindowController : NSWindowController <NSWindowDelegate>
@property NSString *ID;

- (instancetype)initWithID:(NSString *)ID;
@end

@interface TitleBar : NSView
@end

void *Window_New(Window__ w);
void Window_setWebview(NSWindow *win, WKWebView *webview);
void Window_setTitleBar(NSWindow *win, TitleBar *titleBar);

#endif /* window_h */