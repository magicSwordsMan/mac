#include "window.h"
#include "_cgo_export.h"
#include "color.h"

void *Window_New(Window__ w) {
  NSRect contentRect = NSMakeRect(w.X, w.Y, w.Width, w.Height);
  NSUInteger styleMask =
      NSWindowStyleMaskTitled | NSWindowStyleMaskFullSizeContentView |
      NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable |
      NSWindowStyleMaskResizable;

  if (w.Borderless) {
    styleMask = styleMask & NSWindowStyleMaskBorderless;
  }

  if (w.FixedSize) {
    styleMask = styleMask & ~NSWindowStyleMaskResizable;
  }

  if (w.CloseHidden) {
    styleMask = styleMask & ~NSWindowStyleMaskClosable;
  }

  if (w.MinimizeHidden) {
    styleMask = styleMask & ~NSWindowStyleMaskMiniaturizable;
  }

  NSWindow *win = [[NSWindow alloc] initWithContentRect:contentRect
                                              styleMask:styleMask
                                                backing:NSBackingStoreBuffered
                                                  defer:NO];

  if (w.TitlebarHidden) {
    win.titlebarAppearsTransparent = true;
  }

  // Background.
  if (w.Vibrancy != NSVisualEffectMaterialAppearanceBased) {
    NSVisualEffectView *visualEffectView =
        [[NSVisualEffectView alloc] initWithFrame:contentRect];

    visualEffectView.material = w.Vibrancy;
    visualEffectView.blendingMode = NSVisualEffectBlendingModeBehindWindow;
    visualEffectView.state = NSVisualEffectStateActive;
    win.contentView = visualEffectView;
  } else {
    CIColor *backgroundColor = [CIColor colorWithHexString:@"#414244"];
    NSString *bacgroundColorString =
        [NSString stringWithUTF8String:w.BackgroundColor];

    if (bacgroundColorString.length != 0) {
      backgroundColor = [CIColor colorWithHexString:bacgroundColorString];
    }

    win.backgroundColor = [NSColor colorWithCIColor:backgroundColor];
  }

  // Window controller.
  NSString *id = [NSString stringWithUTF8String:w.ID];
  WindowController *controller = [[WindowController alloc] initWithID:id];
  controller.window = win;
  win.windowController = controller;

  // WebView.
  WKWebView *webview =
      Window_NewWebview(controller, [NSString stringWithUTF8String:w.HTML],
                        [NSString stringWithUTF8String:w.ResourcePath]);
  Window_SetWebview(win, webview);
  controller.webview = webview;

  // Titlebar.
  if (w.TitlebarHidden) {
    TitleBar *titleBar = [[TitleBar alloc] init];
    Window_SetTitleBar(win, titleBar);
  }

  [win.windowController showWindow:nil];
  return (__bridge_retained void *)win;
}

WKWebView *Window_NewWebview(WindowController *controller, NSString *HTML,
                             NSString *resourcePath) {
  WKUserContentController *userContentController =
      [[WKUserContentController alloc] init];
  [userContentController addScriptMessageHandler:controller name:@"Call"];

  WKWebViewConfiguration *conf = [[WKWebViewConfiguration alloc] init];
  conf.userContentController = userContentController;

  WKWebView *webView = [[WKWebView alloc] initWithFrame:NSMakeRect(0, 0, 0, 0)
                                          configuration:conf];
  [webView setValue:@(NO) forKey:@"drawsBackground"];
  webView.navigationDelegate = controller;

  // Page loading.
  NSURL *baseURL = [NSURL fileURLWithPath:resourcePath];
  [webView loadHTMLString:HTML baseURL:baseURL];

  while (dispatch_semaphore_wait(controller.sema, DISPATCH_TIME_NOW)) {
    [[NSRunLoop currentRunLoop]
           runMode:NSDefaultRunLoopMode
        beforeDate:[NSDate dateWithTimeIntervalSinceNow:10]];
  }

  return webView;
}

void Window_SetWebview(NSWindow *win, WKWebView *webview) {
  webview.translatesAutoresizingMaskIntoConstraints = NO;
  [win.contentView addSubview:webview];

  [win.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"|[webview]|"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              webview)]];
  [win.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"V:|[webview]|"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              webview)]];
}

void Window_SetTitleBar(NSWindow *win, TitleBar *titleBar) {
  titleBar.translatesAutoresizingMaskIntoConstraints = false;

  [win.contentView addSubview:titleBar];
  [win.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"|[titleBar]|"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              titleBar)]];
  [win.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"V:|[titleBar(==22)]"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              titleBar)]];
}

void Window_CallJS(void *ptr, const char *js) {
  NSWindow *win = (__bridge NSWindow *)ptr;
  WindowController *controller = (WindowController *)win.windowController;

  NSString *javaScript = [NSString stringWithUTF8String:js];
  [controller.webview evaluateJavaScript:javaScript completionHandler:nil];
}

@implementation WindowController
- (instancetype)initWithID:(NSString *)ID {
  self.ID = ID;
  self.sema = dispatch_semaphore_create(0);
  return self;
}

- (void)webView:(WKWebView *)webView
    didFinishNavigation:(WKNavigation *)navigation {
  dispatch_semaphore_signal(self.sema);
}

- (void)userContentController:(WKUserContentController *)userContentController
      didReceiveScriptMessage:(WKScriptMessage *)message {
}
@end

@implementation TitleBar
- (void)mouseDragged:(nonnull NSEvent *)theEvent {
  [self.window performWindowDragWithEvent:theEvent];
}

- (void)mouseUp:(NSEvent *)event {
  NSInteger clickCount = [event clickCount];
  if (2 == clickCount) {
    [self.window zoom:nil];
  }
}
@end