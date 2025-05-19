import React, { forwardRef } from 'react';

import '@components/footer.css';

interface FooterProps extends React.ButtonHTMLAttributes<HTMLDivElement> {
  text: string;
  reload?: boolean;
  settings?: boolean;
  help?: boolean;
  onReloadClick?: () => void;
  onSettingsClick?: () => void;
  onHelpClick?: () => void;
}

const Footer: React.FC<FooterProps> = forwardRef<HTMLDivElement, FooterProps>(
  (
    {
      id,
      text,
      disabled,
      reload = false,
      settings = true,
      help = true,
      onReloadClick,
      onSettingsClick,
      onHelpClick,
      ...props
    },
    ref
  ) => {
    const realId = id || Math.random().toString(36).substring(2, 9);

    return (
      <div id={realId} ref={ref} className="footer-container" {...props}>
        <span className="footer-container__title">
          {text || 'Developed by ONLYOFFICE'}
        </span>
        {reload && (
          <button
            className="footer-container__button"
            onClick={onReloadClick}
            aria-label="Reload"
            type="button"
          >
            <div role="img" className="footer-container__button__reload" />
          </button>
        )}
        {settings && (
          <button
            className="footer-container__button"
            onClick={onSettingsClick}
            aria-label="Settings"
            type="button"
          >
            <div role="img" className="footer-container__button__settings" />
          </button>
        )}
        {help && (
          <button
            className="footer-container__button"
            onClick={onHelpClick}
            aria-label="Help"
            type="button"
          >
            <div role="img" className="footer-container__button__help" />
          </button>
        )}
      </div>
    );
  }
);

Footer.displayName = 'Footer';

export default Footer;
