import React from 'react';

import '@components/select.css';

export interface SelectOption {
  value: string;
  label: string;
}

interface SelectProps {
  options: SelectOption[];
  value: string;
  onChange: (value: string) => void;
  className?: string;
  disabled?: boolean;
}

export const Select: React.FC<SelectProps> = ({
  options,
  value,
  onChange,
  className = '',
  disabled = false,
}) => {
  return (
    <div className="select-container">
      <div className="select-container__wrapper">
        <select
          className={`select-container__field ${className}`}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          disabled={disabled}
        >
          {options.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      </div>
    </div>
  );
};

export default Select;
