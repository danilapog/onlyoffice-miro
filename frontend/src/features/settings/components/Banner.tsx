import React, { forwardRef } from 'react';
import { useTranslation } from 'react-i18next';

import '@features/settings/components/banner.css';

interface BannerProps extends React.HTMLAttributes<HTMLDivElement> {
  url?: string;
}

export const Banner = forwardRef<HTMLDivElement, BannerProps>(
  ({ className, url = 'https://www.onlyoffice.com', ...props }, ref) => {
    const { t } = useTranslation();

    const handleButtonClick = () => {
      window.open(url, '_blank', 'noopener,noreferrer');
    };

    return (
      <div ref={ref} className={`banner ${className || ''}`} {...props}>
        <div className="banner__content">
          <div className="banner__content__info">
            <h3 className="banner__content__info__title">
              {t('features.settings.banner.title')}
            </h3>
            <p className="banner__content__info__description">
              {t('features.settings.banner.description')}
            </p>
          </div>
          <button
            type="button"
            className="banner__content__button"
            onClick={handleButtonClick}
          >
            {t('features.settings.banner.button')}
          </button>
        </div>
      </div>
    );
  }
);

Banner.displayName = 'Banner';

export default Banner;
