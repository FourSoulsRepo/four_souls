import { useCallback, useEffect } from 'react';
import { MainMenu } from './screens/MainMenu';
import { SettingsScreen } from './screens/SettingsScreen';
import { Splash } from './screens/Splash';
import { useApp } from './store';

export function App() {
  const screen = useApp((s) => s.screen);
  const show = useApp((s) => s.show);
  const showMenu = useCallback(() => {
    show('menu');
  }, [show]);

  useEffect(() => {
    void useApp.getState().loadSettings();
  }, []);

  switch (screen) {
    case 'splash':
      return <Splash onDone={showMenu} />;
    case 'menu':
      return <MainMenu />;
    case 'settings':
      return <SettingsScreen />;
  }
}
