import React, { forwardRef } from 'react';

import '@components/input.css';

interface FormInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  fullWidth?: boolean;
  disabled?: boolean;
  helperText?: string;
};

export const FormInput = forwardRef<HTMLInputElement, FormInputProps>(({
  id,
  label,
  error,
  helperText,
  fullWidth = false,
  disabled = false,
  value,
  className = '',
  onChange,
  ...props
}, ref) => {
  const realId = id || Math.random().toString(36).substring(2, 9);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (onChange) onChange(e);
  }

  return (
    <div className="form-input">
      {label && (
        <label
          htmlFor={realId}
          className="form-input__label"
        >
          {label}
        </label>
      )}
      <div className="form-input__container">
        <input
          id={realId}
          ref={ref}
          className={`
            form-input__field
            ${error ? 'form-input__field--error' : ''}
            ${disabled ? 'form-input__field--disabled' : ''}
            ${props.readOnly ? 'form-input__field--readonly' : ''}
            ${className}
          `}
          value={value}
          onChange={handleChange}
          disabled={disabled}
          {...props}
        />
      </div>
      {(error || helperText) && (
        <p className={`form-input__helper-text ${error ? 'form-input__helper-text--error' : 'form-input__helper-text--normal'}`}>
          {error || helperText}
        </p>
      )}
    </div>
  )
});
