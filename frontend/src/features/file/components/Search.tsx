import React, { forwardRef, useState, useEffect, ChangeEvent } from 'react';
import { useFilesStore } from '@features/file/stores/useFileStore';

import './search.css';

interface SearchbarProps extends React.HTMLAttributes<HTMLDivElement> {
}

export const Searchbar = forwardRef<HTMLDivElement, SearchbarProps>(({
  className,
  ...props
}, ref) => {
  const { searchQuery, setSearchQuery } = useFilesStore();
  const [localQuery, setLocalQuery] = useState(searchQuery);

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
      className={`searchbar-container ${className || ''}`}
      ref={ref}
      {...props}
    >
      <div className="searchbar-container__main">
        <div className="searchbar-container__main__icon">
          <img src='/search.svg' />
        </div>
        <input
          className="searchbar-container__main__input"
          type="text"
          placeholder="Search document"
          value={localQuery}
          onChange={handleSearchChange}
        />
        {localQuery && (
          <button className="searchbar-container__main__clear" onClick={handleClearSearch}>
            <img src='/cross.svg' />
          </button>
        )}
      </div>
    </div>
  );
});
