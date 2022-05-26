import simpleRestProvider from 'ra-data-simple-rest';
import { fetchUtils } from 'react-admin';

const httpClient = (url, options = {}) => {
  if (!options.headers) {
    options.headers = new Headers({ Accept: 'application/json' });
  }
  const token = localStorage.getItem('token');
  options.headers.set('Authorization', `Bearer ${token}`);
  return fetchUtils.fetchJson(url, options);
};

const simpleDataProvider = simpleRestProvider('/api', httpClient, 'X-Total-Count')

export const dataProvider = {
  ...simpleDataProvider,
  update: (resource, params) => {
    if (resource !== 'problems' || !params.data.package) {
      // fallback to the default implementation
      return simpleDataProvider.update(resource, params);
    }

    let formData = new FormData();

    formData.append('package', params.data.package.rawFile);

    return httpClient(`/api/${resource}/${params.id}`, {
      method: 'PUT',
      body: formData,
    }).then(({ json }) => ({
      data: { ...params.data, id: json.id },
    }));
  },
  create: (resource, params) => {
    if (resource !== 'problems' || !params.data.package) {
      // fallback to the default implementation
      return simpleDataProvider.create(resource, params);
    }

    let formData = new FormData();

    formData.append('package', params.data.package.rawFile);

    return httpClient(`/api/problems`, {
      method: 'POST',
      body: formData,
    }).then(({ json }) => ({
      data: { ...params.data, id: json.id },
    }));
  }
};;
