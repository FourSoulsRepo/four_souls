import { useCallback, useState } from 'react';
import { MainMenu } from './screens/MainMenu';
import { Splash } from './screens/Splash';

type Screen = 'splash' | 'menu';

export function App() {
  const [screen, setScreen] = useState<Screen>('splash');
  const showMenu = useCallback(() => {
    setScreen('menu');
  }, []);

  return screen === 'splash' ? <Splash onDone={showMenu} /> : <MainMenu />;
}
