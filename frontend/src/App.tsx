import { useCallback, useEffect } from 'react';
import { CardZoom } from './components/CardZoom';
import { CardGallery } from './screens/CardGallery';
import { GameScreen } from './screens/GameScreen';
import { HostScreen } from './screens/HostScreen';
import { JoinScreen } from './screens/JoinScreen';
import { MainMenu } from './screens/MainMenu';
import { SettingsScreen } from './screens/SettingsScreen';
import { Splash } from './screens/Splash';
import { TableScreen } from './screens/TableScreen';
import { useApp, type Screen as ScreenName } from './store';

export function App() {
  const screen = useApp((s) => s.screen);
  const show = useApp((s) => s.show);
  const showMenu = useCallback(() => {
    show('menu');
  }, [show]);

  useEffect(() => {
    void useApp.getState().loadSettings();
    void useApp.getState().loadCards();
  }, []);

  return (
    <>
      <Screen screen={screen} onSplashDone={showMenu} />
      <CardZoom />
    </>
  );
}

function Screen({ screen, onSplashDone }: { screen: ScreenName; onSplashDone: () => void }) {
  switch (screen) {
    case 'splash':
      return <Splash onDone={onSplashDone} />;
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
    case 'cards':
      return <CardGallery />;
  }
}
