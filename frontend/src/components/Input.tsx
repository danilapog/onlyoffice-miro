import React, { forwardRef } from 'react';

import '@components/input.css';

interface FormInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
  fullWidth?: boolean;
  disabled?: boolean;
  required?: boolean;
}

const FormInput = forwardRef<HTMLInputElement, FormInputProps>(
  (
    {
      id,
      label,
      error,
      helperText,
      disabled = false,
      required = false,
      value,
      className = '',
      onChange,
      ...props
    },
    ref
  ) => {
    const realId = id || Math.random().toString(36).substring(2, 9);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      if (onChange) onChange(e);
    };

    return (
      <div className="form-input">
        {label && (
          <label htmlFor={realId} className="form-input__label">
            {label}
            {required && <span className="form-input__label_required">*</span>}
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
            required={required}
            {...props}
          />
        </div>
        {error && (
          <p className="form-input__helper-text form-input__helper-text_error">
            {error}
          </p>
        )}
        {!error && helperText && (
          <p className="form-input__helper-text form-input__helper-text_normal">
            {helperText}
          </p>
        )}
      </div>
    );
  }
);

FormInput.displayName = 'FormInput';

export default FormInput;
