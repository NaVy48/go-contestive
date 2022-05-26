// Taken from https://github.com/marmelab/react-admin/blob/master/packages/ra-ui-materialui/src/auth/LoginForm.tsx
// for mkaing modified version

import Button from '@material-ui/core/Button';
import CardActions from '@material-ui/core/CardActions';
import CircularProgress from '@material-ui/core/CircularProgress';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import PropTypes from 'prop-types';
import { useNotify, useSafeSetState } from 'ra-core';
import React from 'react';
import { fetchUtils } from 'react-admin';
import { Field, Form } from 'react-final-form';
import { Link, useHistory } from 'react-router-dom';

const useStyles = makeStyles(
  theme => ({
    form: {
      padding: '0 1em 1em 1em',
    },
    input: {
      marginTop: '1em',
    },
    button: {
      width: '100%',
    },
    icon: {
      marginRight: theme.spacing(1),
    },
  }),
  { name: 'RaLoginForm' }
);

const Input = ({
  meta: { touched, error }, // eslint-disable-line react/prop-types
  input: inputProps, // eslint-disable-line react/prop-types
  ...props
}) => (
  <TextField
    error={!!(touched && error)}
    helperText={touched && error}
    {...inputProps}
    {...props}
    fullWidth
  />
);

export const RegistrationForm = props => {
  const [loading, setLoading] = useSafeSetState(false);
  const { push } = useHistory();

  const notify = useNotify();
  const classes = useStyles(props);

  const validate = values => {
    const errors = { username: undefined, password: undefined };

    if (!values.email) {
      errors.email = 'Email is required';
    }
    if (!values.username) {
      errors.username = 'Username is required';
    }
    if (!values.password) {
      errors.password = 'Password is required';
    } else if (values.password.length < 8) {
      errors.password = 'Minimum password length is 8';
    }
    return errors;
  };

  const submit = values => {
    setLoading(true);
    fetchUtils
      .fetchJson('/api/auth/register', {
        method: 'POST',
        body: JSON.stringify(values),
      })
      .then(() => {
        setLoading(false);
        push('/login');
      })
      .catch(error => {
        setLoading(false);
        notify(
          typeof error === 'string'
            ? error
            : typeof error === 'undefined' || !error.message
              ? 'Registration error'
              : error.message,
          'warning'
        );
      });
  };

  return (
    <Form
      onSubmit={submit}
      validate={validate}
      render={({ handleSubmit }) => (
        <form onSubmit={handleSubmit} noValidate>
          <div className={classes.form}>
            <div className={classes.input}>
              <Field
                // eslint-disable-next-line jsx-a11y/no-autofocus
                autoFocus
                id="email"
                name="email"
                component={Input}
                label={'email'}
                type="email"
                disabled={loading}
              />
            </div>
            <div className={classes.input}>
              <Field
                id="username"
                name="username"
                component={Input}
                label={'username'}
                disabled={loading}
              />
            </div>
            <div className={classes.input}>
              <Field
                id="password"
                name="password"
                component={Input}
                label={'password'}
                type="password"
                disabled={loading}
                autoComplete="current-password"
              />
            </div>
          </div>
          <CardActions>
            <Button
              variant="contained"
              type="submit"
              color="primary"
              disabled={loading}
              className={classes.button}
            >
              {loading && <CircularProgress className={classes.icon} size={18} thickness={2} />}
              Sign up
            </Button>

            <Button
              type="button"
              component={Link}
              variant="outlined"
              color="primary"
              className={classes.button}
              to="/login"
            >
              Login
            </Button>
          </CardActions>
        </form>
      )}
    />
  );
};

RegistrationForm.propTypes = {
  redirectTo: PropTypes.string,
};
