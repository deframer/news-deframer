import React from 'react';

interface ToggleSwitchProps {
  id: string;
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
}

export const ToggleSwitch: React.FC<ToggleSwitchProps> = ({
  id,
  label,
  checked,
  onChange,
  disabled = false,
}) => {
  const handleToggle = () => {
    if (!disabled) {
      onChange(!checked);
    }
  };

  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        opacity: disabled ? 0.6 : 1,
        cursor: disabled ? 'not-allowed' : 'default',
      }}
    >
      <label htmlFor={id} style={{
        fontSize: '14px',
        fontWeight: 500,
        color: 'var(--text-color)',
        cursor: disabled ? 'not-allowed' : 'pointer',
      }}>
        {label}
      </label>
      <button
        id={id}
        role="switch"
        aria-checked={checked}
        onClick={handleToggle}
        disabled={disabled}
        style={{
          position: 'relative',
          display: 'inline-block',
          width: '44px',
          height: '24px',
          backgroundColor: checked ? 'var(--accent-color)' : 'var(--rating-bg)',
          borderRadius: '34px',
          border: '1px solid var(--btn-border)',
          cursor: disabled ? 'not-allowed' : 'pointer',
          transition: 'background-color 0.2s',
        }}
      >
        <span
          style={{
            position: 'absolute',
            top: '2px',
            left: checked ? '22px' : '2px',
            width: '20px',
            height: '20px',
            borderRadius: '50%',
            backgroundColor: 'white',
            boxShadow: '0 1px 3px rgba(0,0,0,0.2)',
            transition: 'left 0.2s',
          }}
        />
      </button>
    </div>
  );
};
