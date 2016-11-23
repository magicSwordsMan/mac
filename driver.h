#ifndef driver_h
#define driver_h

#import <Cocoa/Cocoa.h>

@interface DriverDelegate : NSObject <NSApplicationDelegate>
@property NSMenu *dock;

- (instancetype)init;
@end

const void *Driver_Init();
void Driver_Run();
void Driver_Terminate();
const char *Driver_Resources();
void Driver_SetAppMenu(const void *menuPtr);
void Driver_SetDockMenu(const void *dockPtr);
void Driver_SetDockIcon(const char *path);
void Driver_SetDockBadge(const char *str);
void Driver_ShowContextMenu(const void *menuPtr);

#endif /* driver_h */