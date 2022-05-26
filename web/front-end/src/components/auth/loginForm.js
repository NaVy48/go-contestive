// Taken from https://github.com/marmelab/react-admin/blob/master/packages/ra-ui-materialui/src/auth/LoginForm.tsx
// for mkaing modified version

import Button from '@material-ui/core/Button';
import CardActions from '@material-ui/core/CardActions';
import CircularProgress from '@material-ui/core/CircularProgress';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import PropTypes from 'prop-types';
import { useLogin, useNotify, useSafeSetState, useTranslate } from 'ra-core';
import React, { FunctionComponent } from 'react';
import { Field, Form } from 'react-final-form';
import { Link } from 'react-router-dom';

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

export const LoginForm = props => {
  const { redirectTo } = props;
  const [loading, setLoading] = useSafeSetState(false);
  const login = useLogin();
  const translate = useTranslate();
  const notify = useNotify();
  const classes = useStyles(props);

  const validate = values => {
    const errors = { username: undefined, password: undefined };

    if (!values.username) {
      errors.username = translate('ra.validation.required');
    }
    if (!values.password) {
      errors.password = translate('ra.validation.required');
    }
    return errors;
  };

  const submit = values => {
    setLoading(true);
    login(values, redirectTo)
      .then(() => {
        setLoading(false);
      })
      .catch(error => {
        setLoading(false);
        notify(
          typeof error === 'string'
            ? error
            : typeof error === 'undefined' || !error.message
              ? 'ra.auth.sign_in_error'
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
                id="username"
                name="username"
                component={Input}
                label={translate('ra.auth.username')}
                disabled={loading}
              />
            </div>
            <div className={classes.input}>
              <Field
                id="password"
                name="password"
                component={Input}
                label={translate('ra.auth.password')}
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
              {translate('ra.auth.sign_in')}
            </Button>
          </CardActions>
        </form>
      )}
    />
  );
};

LoginForm.propTypes = {
  redirectTo: PropTypes.string,
};
