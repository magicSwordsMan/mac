#include "window.h"
#include "_cgo_export.h"
#include "color.h"

const void *Window_New(Window__ w) {
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
  win.delegate = controller;
  win.windowController = controller;
  win.windowController.windowFrameAutosaveName =
      [NSString stringWithUTF8String:w.Title];

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
  return CFBridgingRetain(win);
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

void Window_CallJS(const void *ptr, const char *js) {
  NSWindow *win = (__bridge NSWindow *)ptr;
  WindowController *controller = (WindowController *)win.windowController;

  NSString *javaScript = [NSString stringWithUTF8String:js];
  [controller.webview evaluateJavaScript:javaScript completionHandler:nil];
}

NSRect Window_Frame(const void *ptr) {
  NSWindow *win = (__bridge NSWindow *)ptr;
  return win.frame;
}

void Window_Move(const void *ptr, CGFloat x, CGFloat y) {
  NSWindow *win = (__bridge NSWindow *)ptr;
  CGPoint pos = NSMakePoint(x, y);

  defer([win setFrameOrigin:pos];);
}

void Window_Resize(const void *ptr, CGFloat width, CGFloat height) {
  NSWindow *win = (__bridge NSWindow *)ptr;
  CGRect frame = win.frame;
  frame.size.width = width;
  frame.size.height = height;

  defer([win setFrame:frame display:YES];);
}

void Window_Close(const void *ptr) {
  NSWindow *win = (__bridge NSWindow *)ptr;
  defer([win performClose:nil];);
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

- (void)windowDidMiniaturize:(NSNotification *)notification {
  onWindowMinimize((char *)self.ID.UTF8String);
}

- (void)windowDidDeminiaturize:(NSNotification *)notification {
  onWindowDeminimize((char *)self.ID.UTF8String);
}

- (void)windowDidEnterFullScreen:(NSNotification *)notification {
  onWindowFullScreen((char *)self.ID.UTF8String);
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
  onWindowExitFullScreen((char *)self.ID.UTF8String);
}

- (void)windowDidMove:(NSNotification *)notification {
  onWindowMove((char *)self.ID.UTF8String, self.window.frame.origin.x,
               self.window.frame.origin.y);
}

- (void)windowDidResize:(NSNotification *)notification {
  onWindowResize((char *)self.ID.UTF8String, self.window.frame.size.width,
                 self.window.frame.size.height);
}

- (void)windowDidBecomeKey:(NSNotification *)notification {
  onWindowFocus((char *)self.ID.UTF8String);
}

- (void)windowDidResignKey:(NSNotification *)notification {
  onWindowBlur((char *)self.ID.UTF8String);
}

- (BOOL)windowShouldClose:(id)sender {
  return onWindowClose((char *)self.ID.UTF8String);
}

- (void)windowWillClose:(NSNotification *)notification {
  onWindowCloseFinal((char *)self.ID.UTF8String);
  CFBridgingRelease((__bridge void *)self.window);
  self.window = nil;
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