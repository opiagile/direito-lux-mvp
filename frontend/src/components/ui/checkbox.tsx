import React from 'react';

interface CheckboxProps {
  id?: string;
  checked?: boolean;
  onCheckedChange?: (checked: boolean) => void;
  className?: string;
  disabled?: boolean;
  'aria-label'?: string;
}

export const Checkbox: React.FC<CheckboxProps> = ({
  id,
  checked = false,
  onCheckedChange,
  className = '',
  disabled = false,
  'aria-label': ariaLabel,
}) => {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (onCheckedChange) {
      onCheckedChange(e.target.checked);
    }
  };

  return (
    <input
      id={id}
      type="checkbox"
      checked={checked}
      onChange={handleChange}
      disabled={disabled}
      aria-label={ariaLabel}
      className={`
        w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded 
        focus:ring-blue-500 focus:ring-2 disabled:opacity-50 disabled:cursor-not-allowed
        ${className}
      `}
    />
  );
};

export default Checkbox;