import React, { forwardRef } from 'react';

import '@components/footer.css';

interface FooterProps extends React.ButtonHTMLAttributes<HTMLDivElement> {
  text: string;
  settings?: boolean;
  help?: boolean;
  onSettingsClick?: () => void;
  onHelpClick?: () => void;
}

const Footer: React.FC<FooterProps> = forwardRef<HTMLDivElement, FooterProps>(
  (
    {
      id,
      text,
      disabled,
      settings = true,
      help = true,
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
