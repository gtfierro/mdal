from setuptools import setup, find_packages

print(find_packages())
setup(name='dataclient',
      version='0.0.5',
      description='MDAL client v2',
      packages=find_packages(),
      data_files=[('dataclient', ['dataclient/data.capnp'])],
      include_package_data=True,
      install_requires=[
        'Delorean>=0.6.0',
        'grpcio-tools>=1.14.1',
        'ipython>=6.5.0',
        'msgpack-python>=0.4.2',
        'pandas>=0.23.4',
        'pyaml>=17.12.1',
        'pyarrow>=0.10.0',
        'pycapnp>=0.6.3',
        'rdflib>=4.2.2',
        'requests>=2.19.1',
        'scipy>=1.1.0'
      ],
      zip_safe=False)

