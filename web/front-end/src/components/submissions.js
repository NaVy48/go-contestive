import AssignmentTurnedInIcon from '@material-ui/icons/AssignmentTurnedIn';
import React from 'react';
import {
  Create,
  Datagrid,
  List,
  ReferenceField,
  ReferenceInput,
  SelectInput,
  Show,
  SimpleForm,
  Tab,
  TabbedShowLayout,
  TextField,
  TextInput,
} from 'react-admin';

export const SubmissionList = props => (
  <List {...props} exporter={false} bulkActionButtons={false}>
    <Datagrid rowClick="show">
      <TextField source="id" />
      <ReferenceField source="contestId" reference="contests">
        <TextField source="title" />
      </ReferenceField>
      <ReferenceField source="problemId" reference="problems">
        <TextField source="title" />
      </ReferenceField>

      {props.permission == "admin" && (
        <ReferenceField source="authorId" reference="users">
          <TextField source="username" />
        </ReferenceField>
      )}
      <TextField source="language" />
      <TextField source="status" />
      <TextField source="result" />
    </Datagrid>
  </List>
);

const SrcField = (props) => {
  console.log(props)
  return <pre>{props.record.sourceCode}</pre>
}
export const SubmissionShow = props => (
  <Show {...props}>
    <TabbedShowLayout>
      <Tab label="details">

        <TextField source="id"></TextField>
        <TextField source="createdAt"></TextField>
        <TextField source="updatedAt"></TextField>
        <TextField source="problemId"></TextField>
        <TextField source="problemRevId"></TextField>
        <TextField source="contestId"></TextField>
        <TextField source="authorId"></TextField>
        <TextField source="language"></TextField>
        <TextField source="result"></TextField>
      </Tab>
      <Tab label="Source Code">
        <SrcField></SrcField>
      </Tab>
    </TabbedShowLayout>
  </Show>
);

export const SubmissionCreate = ({ contestId, ...props }) => (
  <Create {...props}>
    <SimpleForm defaultValue={{ contestId }}>
      <ReferenceInput source="problemId" reference="problems" filter={{ "cp.contestid": contestId }}>
        <SelectInput optionText="title" />
      </ReferenceInput>
      <SelectInput
        source="language"
        label="Language"
        choices={[
          { id: 'CPP11', name: 'CPP11' },
          // { id: 'JAVA', name: 'JAVA' },
        ]}
      />
      <TextInput fullWidth multiline source="sourceCode" />
    </SimpleForm>
  </Create>
);



export const SubmissionsComponents = {
  list: SubmissionList,
  create: SubmissionCreate,
  show: SubmissionShow,
  icon: AssignmentTurnedInIcon,
};
