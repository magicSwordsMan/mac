#include "driver.h"
#include "_cgo_export.h"

@implementation DriverDelegate
- (instancetype)init {
  self.dock = [[NSMenu alloc] initWithTitle:@""];
  return self;
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
  onLaunch();
}

- (void)applicationDidBecomeActive:(NSNotification *)aNotification {
  onFocus();
}

- (void)applicationDidResignActive:(NSNotification *)aNotification {
  onBlur();
}

- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender
                    hasVisibleWindows:(BOOL)flag {
  onReopen(flag);
  return YES;
}

- (BOOL)application:(NSApplication *)theApplication
           openFile:(NSString *)filename {
  onFileOpen((char *)[filename UTF8String]);
  return YES;
}

- (NSApplicationTerminateReply)applicationShouldTerminate:
    (NSApplication *)sender {
  return onTerminate();
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
  onFinalize();
}
@end

void *Driver_Init() {
  [NSApplication sharedApplication];
  [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
  [NSApp activateIgnoringOtherApps:YES];

  DriverDelegate *delegate = [[DriverDelegate alloc] init];
  NSApp.delegate = delegate;

  return (void *)CFBridgingRetain(delegate);
}

void Driver_Run() { [NSApp run]; }

void Driver_Terminate() { [NSApp terminate:NSApp]; }