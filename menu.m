#include "menu.h"
#include "_cgo_export.h"
#include "mac.h"

const void *Menu_New(Menu__ m) {
  Menu *menu = [[Menu alloc] init];
  menu.ID = [NSString stringWithUTF8String:m.ID];

  return CFBridgingRetain(menu);
}

void Menu_Mount(const void *ptr, const char *rootID) {
  Menu *menu = (__bridge Menu *)ptr;
  MenuContainer *container =
      [menu.Elems objectForKey:[NSString stringWithUTF8String:rootID]];
  menu.Root = container;
  container.delegate = menu;
}

void Menu_Dismount(const void *ptr) {
  Menu *menu = (__bridge Menu *)ptr;
  [menu dismountElement:menu.Root];
}

void Menu_MountContainer(const void *ptr, MenuContainer__ c) {
  Menu *menu = (__bridge Menu *)ptr;
  NSString *containerID = [NSString stringWithUTF8String:c.ID];

  MenuContainer *container = [menu.Elems objectForKey:containerID];
  if (container == nil) {
    container = [[MenuContainer alloc]
        initWithTitle:[NSString stringWithUTF8String:c.Label]];
    container.ID = containerID;
    [menu.Elems setObject:container forKey:containerID];
    return;
  }

  for (int i = 0; i < container.numberOfItems; i++) {
    [menu dismountElement:container.itemArray[i]];
  }

  defer([container removeAllItems];);
}

void Menu_MountItem(const void *ptr, MenuItem__ it) {
  Menu *menu = (__bridge Menu *)ptr;
  NSString *itemID = [NSString stringWithUTF8String:it.ID];

  MenuItem *item = [menu.Elems objectForKey:itemID];
  if (item == nil) {

    item = [[MenuItem alloc] init];
    item.ID = itemID;

    [menu.Elems setObject:item forKey:itemID];
  }

  item.title = [NSString stringWithUTF8String:it.Label];
  item.OnClick = [NSString stringWithUTF8String:it.OnClick];
  item.enabled = !it.Disabled;
  item.IsSeparator = it.Separator;

  NSString *icon = [NSString stringWithUTF8String:it.Icon];
  if (icon.length != 0) {
    item.image = [[NSImage alloc] initByReferencingFile:icon];
  } else {
    item.image = nil;
  }

  [item setSelector:[NSString stringWithUTF8String:it.Selector]];
  [item setShortcut:[NSString stringWithUTF8String:it.Shortcut]];
  [item setSeparator];
}

void Menu_Associate(const void *ptr, const char *parentID,
                    const char *childID) {
  Menu *menu = (__bridge Menu *)ptr;
  MenuContainer *parent =
      [menu.Elems objectForKey:[NSString stringWithUTF8String:parentID]];
  id child = [menu.Elems objectForKey:[NSString stringWithUTF8String:childID]];

  //  child is a MenuItem.
  if ([child isKindOfClass:[MenuItem class]]) {
    MenuItem *c = (MenuItem *)child;

    defer([parent addItem:c];[c setSeparator];);
    return;
  }

  //  child is a MenuContainer.
  MenuContainer *container = (MenuContainer *)child;
  MenuItem *item = [[MenuItem alloc] init];
  item.title = container.title;
  item.submenu = container;
  defer([parent addItem:item];);
}

void Menu_Clear(const void *ptr) {
  Menu *menu = (__bridge Menu *)ptr;

  menu.Root = nil;
  menu.Elems = [NSMutableDictionary dictionary];
}

@implementation MenuItem
- (void)setSelector:(NSString *)selector {
  if (!self.enabled) {
    self.action = nil;
    return;
  }

  if (self.hasSubmenu) {
    self.action = @selector(submenuAction:);
    return;
  }

  SEL action = @selector(clicked:);
  self.target = self;

  if (selector.length > 0) {
    self.target = nil;
    action = NSSelectorFromString(selector);
  }
  self.action = action;
}

- (void)setShortcut:(NSString *)shortcut {
  if (shortcut.length == 0) {
    return;
  }

  NSArray *keys = [shortcut componentsSeparatedByString:@"+"];
  self.keyEquivalentModifierMask = 0;

  for (NSString *k in keys) {
    if ([k isEqual:@"meta"]) {
      self.keyEquivalentModifierMask |= NSEventModifierFlagCommand;
    } else if ([k isEqual:@"ctrl"]) {
      self.keyEquivalentModifierMask |= NSEventModifierFlagControl;
    } else if ([k isEqual:@"alt"]) {
      self.keyEquivalentModifierMask |= NSEventModifierFlagOption;
    } else if ([k isEqual:@"shift"]) {
      self.keyEquivalentModifierMask |= NSEventModifierFlagShift;
    } else if ([k isEqual:@"fn"]) {
      self.keyEquivalentModifierMask |= NSEventModifierFlagFunction;
    } else if ([k isEqual:@""]) {
      self.keyEquivalent = @"+";
    } else {
      self.keyEquivalent = k;
    }
  }
}

- (void)setSeparator {
  NSMenu *container = self.menu;
  if (container == nil) {
    return;
  }

  if (self.IsSeparator && self.SeparatorItem == nil) {
    NSInteger idx = [container indexOfItem:self];

    self.SeparatorItem = [NSMenuItem separatorItem];
    [container insertItem:self.SeparatorItem atIndex:idx + 1];
    return;
  }

  if (!self.IsSeparator && self.SeparatorItem != nil) {
    [container removeItem:self.SeparatorItem];
    self.SeparatorItem = nil;
    return;
  }
}

- (void)clicked:(id)sender {
  onMenuItemClick((char *)self.ID.UTF8String, (char *)self.OnClick.UTF8String);
}
@end

@implementation MenuContainer
@end

@implementation Menu
- (instancetype)init {
  self.Elems = [NSMutableDictionary dictionary];
  return self;
}

- (void)dismountElement:(id)elem {
  //  elem is a MenuContainer.
  if ([elem isKindOfClass:[MenuContainer class]]) {
    MenuContainer *container = (MenuContainer *)elem;

    for (int i = 0; i < container.numberOfItems; i++) {
      [self dismountElement:container.itemArray[i]];
    }

    if ([self.Elems objectForKey:container.ID] != nil) {
      [self.Elems removeObjectForKey:container.ID];
    }

    return;
  }

  //  elem is a MenuItem.
  if ([elem isKindOfClass:[MenuItem class]]) {
    MenuItem *item = (MenuItem *)elem;

    if (item.submenu != nil) {
      [self dismountElement:item.submenu];
    }

    if ([self.Elems objectForKey:item.ID] != nil) {
      [self.Elems removeObjectForKey:item.ID];
    }
  }
}

- (void)menuDidClose:(NSMenu *)menu {
  onMenuCloseFinal((char *)self.ID.UTF8String);
  CFBridgingRelease((__bridge void *)self);
}
@end