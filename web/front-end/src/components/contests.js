import StarIcon from '@material-ui/icons/Star';
import React from 'react';
import {
  Create,
  Datagrid,
  DateField,
  DateTimeInput,
  Edit,
  FormTab,
  List,
  NumberField,
  ReferenceArrayField,
  ReferenceArrayInput,
  ReferenceField,
  SelectArrayInput,
  Show,
  SimpleForm,
  Tab,
  TabbedForm,
  TabbedShowLayout,
  TextField,
  TextInput,
} from 'react-admin';

import { SubmissionCreate, SubmissionList } from './submissions';
import styled from 'styled-components'

export const ContestList = props => (
  <List {...props} exporter={false} bulkActionButtons={false}>
    <Datagrid rowClick="show">

      {props.permission == "admin" && (
        <>
          <TextField source="id" />
          <ReferenceField source="authorId" reference="users">
            <TextField source="username" />
          </ReferenceField>
        </>
      )}
      <TextField source="title" />
      <DateField source="startTime" showTime />
      <DateField source="endTime" showTime />
    </Datagrid>
  </List>
);

export const ContestEdit = props => (
  <Edit {...props} undoable={false}>
    <TabbedForm>
      <FormTab label="summary">
        <TextInput source="title" />
        <DateTimeInput source="startTime" />
        <DateTimeInput source="endTime" />
      </FormTab>
      <FormTab label="Problems" path="problems">
        <ReferenceArrayInput
          source="problems"
          reference="problems"
          fullWidth
          label="Problems"
        >
          <SelectArrayInput optionText="name" fullWidth />
        </ReferenceArrayInput>
      </FormTab>
      <FormTab label="Users" path="user">
        <ReferenceArrayInput
          source="users"
          reference="users"
          fullWidth
          label="Users"
        >
          <SelectArrayInput optionText="username" fullWidth />
        </ReferenceArrayInput>
      </FormTab>
    </TabbedForm>
  </Edit>
);

const Submit = ({ id }) => {
  const CreateProps = {
    basePath: '/submissions',
    resource: 'submissions',
  };
  return <SubmissionCreate contestId={id} {...CreateProps} />;
};
const Submitions = ({ permission, id }) => {
  const ShowProps = {
    basePath: '/submissions',
    hasShow: true,
    resource: 'submissions',
  };
  return (
    <List {...ShowProps} filterDefaultValues={{ contestId: id }} exporter={false} bulkActionButtons={false} actions={false}>
      <Datagrid rowClick="show">
        <ReferenceField source="problemId" reference="problems">
          <TextField source="title" />
        </ReferenceField>
        {permission == "admin" && (
          <ReferenceField source="authorId" reference="users">
            <TextField source="username" />
          </ReferenceField>
        )}
        <TextField source="status" />
        <TextField source="result" />
      </Datagrid>
    </List>
  )
};

const FullWidth = styled.div`
 .MuiFormControl-root {
   display: flex
 }
`

export const ContestShow = props => {
  console.log('ContestShow');
  console.log(props);
  return (
    <Show {...props}>
      <TabbedShowLayout>
        <Tab label="summary">

          {props.permission == "admin" && (
            <>
              <TextField source="id" />
              <ReferenceField source="authorId" reference="users">
                <TextField source="username" />
              </ReferenceField>
            </>
          )}
          <TextField source="title" />
          <DateField source="startTime" showTime />
          <DateField source="endTime" showTime />
        </Tab>
        <Tab label="Problems" path="tests">
          <FullWidth>
            <ReferenceArrayField
              source="problems"
              reference="problems"
              fullWidth
              sort={{ field: 'title', order: 'ASC' }}
            >
              <Datagrid rowClick="show"  >
                <TextField source="title" />
                <NumberField source="memoryLimit" />
                <NumberField source="timeLimit" />
              </Datagrid>
            </ReferenceArrayField>
          </FullWidth>
        </Tab>
        <Tab label="Submit" path="submit">
          <Submit id={props.id} />
        </Tab>
        <Tab label="Submitions" path="submitions">
          <Submitions id={props.id} />
        </Tab>
      </TabbedShowLayout>
    </Show>
  );
};

export const ContestCreate = props => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="title" />
      <DateTimeInput source="startTime" />
      <DateTimeInput source="endTime" />
    </SimpleForm>
  </Create>
);
export const ContestsComponents = {
  list: ContestList,
  create: ContestCreate,
  show: ContestShow,
  edit: ContestEdit,
  icon: StarIcon,
};

export const UserContestsComponents = {
  list: ContestList,
  show: ContestShow,
  icon: StarIcon,
};
