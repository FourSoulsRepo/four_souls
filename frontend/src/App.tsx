import { useCallback, useEffect, useRef } from 'react';
import { CardZoom } from './components/CardZoom';
import { NoticeOverlay } from './components/Overlays';
import { Effects } from './components/Effects';
import { Sky } from './components/Sky';
import { CardGallery } from './screens/CardGallery';
import { GameScreen } from './screens/GameScreen';
import { HostScreen } from './screens/HostScreen';
import { JoinScreen } from './screens/JoinScreen';
import { MainMenu } from './screens/MainMenu';
import { SettingsScreen } from './screens/SettingsScreen';
import { Splash } from './screens/Splash';
import { TableScreen } from './screens/TableScreen';
import { WelcomeScreen } from './screens/WelcomeScreen';
import { onSky } from './sky';
import { useApp, type Screen as ScreenName } from './store';

export function App() {
  const screen = useApp((s) => s.screen);
  const show = useApp((s) => s.show);
  const loadingRef = useRef<Promise<void> | null>(null);
  // After the splash: the menu, or the nickname first on the first start
  // (ST-06). The settings may still be loading.
  const showMenu = useCallback(() => {
    void (loadingRef.current ?? Promise.resolve()).then(() => {
      show(useApp.getState().settings.nickname ? 'menu' : 'welcome');
    });
  }, [show]);

  useEffect(() => {
    loadingRef.current = useApp.getState().loadSettings();
    void useApp.getState().loadCards();
  }, []);

  return (
    <>
      {onSky(screen) && <Sky />}
      <Screen screen={screen} onSplashDone={showMenu} />
      <Effects />
      <CardZoom />
      <NoticeOverlay />
    </>
  );
}

function Screen({ screen, onSplashDone }: { screen: ScreenName; onSplashDone: () => void }) {
  switch (screen) {
    case 'splash':
      return <Splash onDone={onSplashDone} />;
    case 'welcome':
      return <WelcomeScreen />;
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
