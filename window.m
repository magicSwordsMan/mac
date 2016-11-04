#include "window.h"
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
  NSString *id = [NSString stringWithUTF8String:w.BackgroundColor];
  WindowController *winController = [[WindowController alloc] initWithID:id];
  winController.window = win;
  win.windowController = winController;

  // WebView.
  WKWebView *webView = [[WKWebView alloc] initWithFrame:NSMakeRect(0, 0, 0, 0)];
  Window_setWebview(win, webView);

  // Titlebar.
  TitleBar *titleBar = [[TitleBar alloc] init];
  Window_setTitleBar(win, titleBar);

  [win.windowController showWindow:nil];
  return (__bridge_retained void *)win;
}

void Window_setWebview(NSWindow *win, WKWebView *webView) {
  webView.translatesAutoresizingMaskIntoConstraints = NO;
  [win.contentView addSubview:webView];

  [win.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"|[webView]|"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              webView)]];
  [win.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"V:|[webView]|"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              webView)]];
}

void Window_setTitleBar(NSWindow *win, TitleBar *titleBar) {
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

@implementation WindowController
- (instancetype)initWithID:(NSString *)ID {
  self.ID = ID;
  return self;
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