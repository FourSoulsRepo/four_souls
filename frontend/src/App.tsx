import { useCallback, useEffect } from 'react';
import { GameScreen } from './screens/GameScreen';
import { HostScreen } from './screens/HostScreen';
import { JoinScreen } from './screens/JoinScreen';
import { MainMenu } from './screens/MainMenu';
import { SettingsScreen } from './screens/SettingsScreen';
import { Splash } from './screens/Splash';
import { TableScreen } from './screens/TableScreen';
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
    case 'host':
      return <HostScreen />;
    case 'join':
      return <JoinScreen />;
    case 'table':
      return <TableScreen />;
    case 'game':
      return <GameScreen />;
  }
}
