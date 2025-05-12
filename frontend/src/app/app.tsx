import React, { useEffect, useRef, useState } from 'react';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { HashRouter, Routes, Route, useLocation } from 'react-router-dom';
import { CSSTransition, SwitchTransition } from 'react-transition-group';

import { useApplicationStore } from '@stores/useApplicationStore';
import { ManagerPage } from '@app/pages/manager';
import { CreationPage } from '@app/pages/creation';
import { SettingsPage } from '@app/pages/settings';
import { useFilesStore } from '@features/file/stores/useFileStore';
import { Button } from '@components/Button';

import '@app/transitions.css';
import '@i18n/config';

const App = () => {
  const { loading, authorized, admin, reload, refresh } = useApplicationStore();
  const { refreshDocuments } = useFilesStore();
  const location = useLocation();
  const nodeRef = useRef(null);
  const [prevPathname, setPrevPathname] = useState(location.pathname);
  
  useEffect(() => {
    reload();
    refreshDocuments();
  }, []);

  useEffect(() => {
    const isBack = prevPathname.length > location.pathname.length;
    if (!location.state) {
      location.state = { isBack };
    }
    setPrevPathname(location.pathname);
  }, [location.pathname]);

  if (loading) return <div>Loading...</div>

  if (!authorized) return (
    <>
      <Button name='refresh' onClick={refresh} />
      <div>Please install/reinstall the app</div>
    </>
  );

  return (
    <div className="page-container">
      <SwitchTransition mode="out-in">
        <CSSTransition
          key={location.pathname}
          nodeRef={nodeRef}
          timeout={300}
          classNames={location.state?.isBack ? "page-back" : "page-forward"}
          unmountOnExit
        >
          <div ref={nodeRef}>
            <Routes location={location}>
              <Route path="/" element={<ManagerPage />} />
              <Route path="/create" element={<CreationPage />} />
              {admin && <Route path="/settings" element={<SettingsPage />} />}
            </Routes>
          </div>
        </CSSTransition>
      </SwitchTransition>
    </div>
  );
};

const rootElement = document.getElementById('root');
if (rootElement) {
  createRoot(rootElement).render(
    <StrictMode>
      <HashRouter>
        <App />
      </HashRouter>
    </StrictMode>
  );
}
