import React, { forwardRef, useState, useEffect, ChangeEvent } from 'react';

import useFilesStore from '@features/file/stores/useFileStore';

import '@features/file/components/search.css';

interface SearchbarProps extends React.HTMLAttributes<HTMLDivElement> {}

export const Searchbar = forwardRef<HTMLDivElement, SearchbarProps>(
  ({ className, ...props }, ref) => {
    const { searchQuery, setSearchQuery, initialized, loading } =
      useFilesStore();
    const [localQuery, setLocalQuery] = useState(searchQuery);

    const disabled = loading && !initialized;

    useEffect(() => {
      const timer = setTimeout(() => {
        setSearchQuery(localQuery);
      }, 300);

      return () => clearTimeout(timer);
    }, [localQuery, setSearchQuery]);

    const handleSearchChange = (e: ChangeEvent<HTMLInputElement>) => {
      setLocalQuery(e.target.value);
    };

    const handleClearSearch = () => {
      setLocalQuery('');
      setSearchQuery('');
    };

    return (
      <div
        className={`searchbar-container ${className || ''} ${disabled ? 'searchbar-container__disabled' : ''}`}
        ref={ref}
        {...props}
      >
        <div className="searchbar-container__main">
          <div className="searchbar-container__main__icon">
            <img src="/search.svg" alt="Search icon" />
          </div>
          <input
            className="searchbar-container__main__input"
            type="text"
            placeholder="Search document"
            value={localQuery}
            onChange={handleSearchChange}
            disabled={disabled}
          />
          {localQuery && (
            <button
              type="button"
              className="searchbar-container__main__clear"
              onClick={handleClearSearch}
              disabled={disabled}
              aria-label="Clear search"
            >
              <img src="/cross.svg" alt="Clear search" />
            </button>
          )}
        </div>
      </div>
    );
  }
);

Searchbar.displayName = 'Searchbar';

export default Searchbar;
